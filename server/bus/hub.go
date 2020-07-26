package bus

import (
	"fmt"

	"github.com/google/uuid"
	//	"github.com/orcaman/concurrent-map"
)

// Hub for handling client connections on channels
type Hub struct {
	// Registered clients.
	channels map[uuid.UUID]*Channel

	// Unregistered clients
	clientPool map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *HubMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type HubMessage struct {
	data   []byte
	client *Client
}

type Channel struct {
	clients map[uuid.UUID]*Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *HubMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		channels:   make(map[uuid.UUID]*Channel),
		clientPool: make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.RegisterClient(client)
		case client := <-h.unregister:
			h.UnregisterClient(client)
		case message := <-h.broadcast:
			h.Broadcast(message)
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	if client.channelID != uuid.Nil {
		channel, ok := h.channels[client.channelID]
		if !ok {
			channel = &Channel{
				clients: make(map[uuid.UUID]*Client),
			}
			h.channels[client.channelID] = channel
		}
		channel.clients[client.ID] = client
		delete(h.clientPool, client)

		fmt.Printf("Registered %s to %s\n", client.ID, client.channelID)
	}
}

func (h *Hub) UnregisterClient(client *Client) {
	if channel, ok := h.channels[client.channelID]; ok {
		delete(channel.clients, client.ID)
		close(client.send)

		if len(channel.clients) == 0 {
			delete(h.channels, client.channelID)
		}
	}
}

func (h *Hub) Broadcast(message *HubMessage) {
	from := message.client.ID
	channelID := message.client.channelID

	if channel, ok := h.channels[channelID]; ok {
		fmt.Printf("Received message from %s to broadcast to %d users on channel %s\n", message.client.ID, len(channel.clients)-1, channelID)

		for to, toClient := range channel.clients {
			// send to everyone else in the channel
			if to == from {
				continue
			}

			fmt.Printf("Sending to %s in channel\n", to)
			select {
			case toClient.send <- message.data:
			default:
				close(toClient.send)
				delete(channel.clients, to)
			}
		}
	}
}
