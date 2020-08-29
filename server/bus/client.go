package bus

import (
	"encoding/json"
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
		// TODO: log error here
		return
	}

	if !config.Get().IsValidClient(clientVersion) {
		fmt.Println("Rejecting client", id)
		c.Unregister()
		return
	}

	c.ID = id
	c.channelID = channelID
	c.hub.register <- c
}

func (c *Client) Unregister() {
	fmt.Println("Unregistering client", c.ID)
	close(c.send)
	c.conn.Close()
}

// ReadPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		fmt.Println("Closing read pump for", c.ID)
		c.Unregister()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message from connection for client %s: %v\n", c.ID, err)
			return
		}

		msg, err := actions.NewMessageFrom(message)
		if err != nil {
			fmt.Printf("Error parsing message into Message struct: %+v\n%s\n\n", err, message)
			continue
		}

		switch msg.Type {
		// register is a special case that we don't want to rebroadcast
		case actions.Register:
			c.Register(msg.From, msg.Channel, msg.ClientVersion)
			continue

		case actions.SetGameID:
			if ok := c.hub.SetGameIdForChannel(c.channelID, msg.GameID); !ok {
				fmt.Println("Dropping client", c.ID)
				return
			}
			continue

		case actions.ApplyProgress:
			apMsg, err := c.hub.ApplyProgressToChannel(msg.Payload, c.channelID)
			if err != nil {
				fmt.Printf("Error applying progress to channel %+v\n", err)
				continue
			}

			data, err := json.Marshal(&apMsg)
			if err != nil {
				fmt.Printf("Error marshalling latest message into byte array %+v\n", err)
				continue
			}

			raw := json.RawMessage(data)
			msg.Payload = &raw
			fmt.Printf("Updated apply progress message -- %+v\n %+v\n", msg, apMsg)
		}

		if actions.IsIgnoredAction(msg.Type) {
			fmt.Println("Ignoring", msg.Type)
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
		fmt.Println("Closing write pump for", c.ID)
		ticker.Stop()
		c.hub.unregister <- c
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// The hub closed the channel.
			if !ok {
				fmt.Printf("client %s channel closed\n", c.ID)
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
