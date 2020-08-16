package bus

import (
	"fmt"
	"time"

	"github.com/salmaanrizvi/crossword-party/server/actions"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"
)

// Hub for handling client connections on channels
type Hub struct {
	// Registered clients.
	channels cmap.ConcurrentMap //map[uuid.UUID]*Channel

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
	action actions.Action
}

type Channel struct {
	clients cmap.ConcurrentMap //map[uuid.UUID]*Client
	gameID  int
	// register   chan *Client
	// unregister chan *Client
}

// NewHub ...
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *HubMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		channels:   cmap.New(), //make(map[uuid.UUID]*Channel),
	}
}

// Run ...
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

func (h *Hub) Stats(timeInSeconds int) {
	ticker := time.NewTicker(time.Duration(timeInSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				return
			}

			h.printStats()
		}
	}
}

func (h *Hub) printStats() {
	fmt.Println(time.Now())
	fmt.Printf("Channels Count: %d\n", h.channels.Count())
	for channelID, val := range h.channels.Items() {
		if channel, ok := val.(*Channel); ok {
			fmt.Printf("Channel %s: %d connected clients\n", channelID, channel.clients.Count())
		}
	}
}

// RegisterClient ...
func (h *Hub) RegisterClient(client *Client) {
	if client.channelID == uuid.Nil {
		return
	}

	var channel *Channel
	tmp, ok := h.channels.Get(client.channelID.String())
	if !ok {
		channel = &Channel{clients: cmap.New()}
		h.channels.Set(client.channelID.String(), channel)
	} else {
		channel = tmp.(*Channel)
	}

	channel.clients.Set(client.ID.String(), client)
	fmt.Printf("Registered client %s to channel %s\n", client.ID, client.channelID)
}

// UnregisterClient ...
func (h *Hub) UnregisterClient(client *Client) {
	if client.channelID == uuid.Nil {
		fmt.Printf("Cant unregister client (%s) from nil channel uuid\n", client.ID)
		return
	}

	tmp, ok := h.channels.Get(client.channelID.String())
	if !ok {
		fmt.Printf("Couldnt find channel %s to unregister client from\n", client.channelID)
		return
	}

	channel, ok := tmp.(*Channel)
	if !ok {
		// TODO - clean up this in our map
		fmt.Printf("Found non-channel reference in map, we should clean this up...")
		return
	}

	// Remove client from channel & close the connection
	channel.clients.Remove(client.ID.String())

	// optionally, remove channel from hub if it has no clients left
	if channel.clients.Count() == 0 {
		h.channels.Remove(client.channelID.String())
		fmt.Printf("Channel %s is empty - removing\n", client.channelID)
	}
}

func (h *Hub) SetGameIdForChannel(channelID uuid.UUID, gameID int) bool {
	tmp, ok := h.channels.Get(channelID.String())
	if !ok {
		fmt.Printf("Couldnt set game id for channel: %s\n", channelID)
		return false
	}

	channel, ok := tmp.(*Channel)
	if !ok {
		fmt.Printf("Couldnt cast channel with id %s to channel type to set game id\n", channelID)
		return false
	}

	if channel.gameID != 0 && channel.gameID != gameID {
		fmt.Printf("Game %d is already in session in channel %s. Can not set it to %d\n", channel.gameID, channelID, gameID)
		return false
	}

	channel.gameID = gameID
	fmt.Printf("Set channel %s to game id %d\n", channelID, gameID)
	return true
}

// Broadcast ...
func (h *Hub) Broadcast(message *HubMessage) {
	from := message.client.ID.String()
	channelID := message.client.channelID.String()

	tmp, ok := h.channels.Get(channelID)
	if !ok {
		fmt.Printf("Couldnt find channel %s to broadcast message to\n", channelID)
		return
	}

	channel, ok := tmp.(*Channel)
	if !ok {
		fmt.Printf("Channel was not of type *Channel -- actually was %T", tmp)
		return
	}

	fmt.Printf("Broadcasting %s\n", message.action)
	for to, _client := range channel.clients.Items() {
		// send to everyone else in the channel
		if to == from {
			continue
		}

		client, ok := _client.(*Client)
		if !ok {
			fmt.Printf("Unknown client type (%T) to broadcast to... skipping\n", _client)
			continue
		}

		select {
		case client.send <- message.data:
		default:
			h.UnregisterClient(client)
		}

		fmt.Printf("Sent to %s in channel\n", to)
	}
}
