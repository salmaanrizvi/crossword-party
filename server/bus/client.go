package bus

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/salmaanrizvi/crossword-party/actions"
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
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) Register(id, channelID uuid.UUID) {
	if id != uuid.Nil {
		c.ID = id
	}

	if channelID != uuid.Nil {
		c.channelID = channelID
	}

	c.hub.register <- c
}

func (c *Client) Unregister() {
	c.hub.unregister <- c
	c.conn.Close()
}

// ReadPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Unregister()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg, err := actions.NewMessageFrom(message)
		if err != nil {
			fmt.Printf("Error parsing message into Message struct: %+v\n%s\n\n", err, message)
			continue
		}

		fmt.Printf("Parsed message %+v\n", msg)

		switch msg.Type {
		case actions.Register.String():
			c.Register(msg.From, msg.Channel)
			continue

		case actions.ApplyProgress.String():
			_, err := actions.NewApplyProgessMessageFrom(message)
			if err != nil {
				fmt.Printf("Error parsing payload into ApplyProgressMessage struct: %+v\n%s\n\n", err, message)
				// continue
			}
			// fmt.Printf("Successfully parsed ApplyProgressMessage: %+v", ap)
			// only explicitly let action types we care about pass through
			// default:
			// 	continue
		}

		c.hub.broadcast <- &HubMessage{data: message, client: c}
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
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				fmt.Printf("closing connection because msg error: %+v\n", message)
				// The hub closed the channel.
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
