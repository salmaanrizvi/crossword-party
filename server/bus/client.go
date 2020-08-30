package bus

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/salmaanrizvi/crossword-party/server/actions"
	"github.com/salmaanrizvi/crossword-party/server/config"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	channelID uuid.UUID
	ID        uuid.UUID
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	// TODO: SetCloseHandler, SetPingHandler, SetPongHandler
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) Register(id, channelID uuid.UUID, clientVersion string) {
	if id == uuid.Nil || channelID == uuid.Nil {
		config.Logger().Debugf(
			"Rejecting client",
			"client_id", id,
			"channel_id", channelID,
			"client_version", clientVersion,
			"reason", "Nil client or channel id provided",
		)
		return
	}

	if !config.Get().IsValidClient(clientVersion) {
		c.Unregister(fmt.Errorf("Client %s had unsupported version %s", id, clientVersion))
		return
	}

	c.ID = id
	c.channelID = channelID
	c.hub.register <- c
}

func (c *Client) Unregister(e error) {
	config.Logger().Infow(
		"Unregistering client",
		"client_id", c.ID,
		"channel_id", c.channelID,
		"reason", e.Error(),
	)
	close(c.send)
	c.conn.Close()
}

// ReadPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	var e error
	defer func() {
		// config.Logger().Debugw("Closing read pump for", "client_id", c.ID)
		c.Unregister(e)
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			e = fmt.Errorf("Error reading message from connection: %s", err.Error())
			return
		}

		msg, err := actions.NewMessageFrom(message)
		if err != nil {
			config.Logger().Errorw("Error parsing message into struct: %s \n\n%s", err.Error(), message)
			continue
		}

		switch msg.Type {
		// register is a special case that we don't want to rebroadcast
		case actions.Register:
			c.Register(msg.From, msg.Channel, msg.ClientVersion)
			continue

		case actions.SetGameID:
			if e = c.hub.SetGameIdForChannel(c.channelID, msg.GameID); e != nil {
				return
			}

			config.Logger().Infow(
				"Set game id for channel",
				"game_id", msg.GameID,
				"channel_id", c.channelID,
				"client_id", c.ID,
			)
			continue

		case actions.ApplyProgress:
			apMsg, err := c.hub.ApplyProgressToChannel(msg.Payload, c.channelID)
			if err != nil {
				config.Logger().Errorw("Error applying progress from client %s to channel %s: %s", c.ID, c.channelID, err.Error())
				continue
			}

			msg.Payload = apMsg
			config.Logger().Debugw("Updated apply progress message", "channel_id", c.channelID)
		}

		if actions.IsIgnoredAction(msg.Type) {
			config.Logger().Debugw("Ignoring message", "type", msg.Type)
			continue
		}

		c.hub.broadcast <- &HubMessage{data: message, client: c, action: msg.Type, sendAll: msg.Type == actions.ApplyProgress}
	}
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregister <- c
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// The hub closed the channel.
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
