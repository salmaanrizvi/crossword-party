package bus

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/salmaanrizvi/crossword-party/server/actions"
	"github.com/salmaanrizvi/crossword-party/server/config"

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
	data    []byte
	client  *Client
	action  actions.Action
	sendAll bool
}

type Channel struct {
	clients  cmap.ConcurrentMap //map[uuid.UUID]*Client
	gameID   int
	Progress *actions.ApplyProgressMessage
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

func (h *Hub) GetChannel(channelID uuid.UUID) (*Channel, error) {
	tmp, ok := h.channels.Get(channelID.String())
	if !ok {
		return nil, fmt.Errorf("Channel %s not found in list of channels", channelID)
	}

	channel, ok := tmp.(*Channel)
	if !ok {
		// TODO - clean up this in our map
		return nil, fmt.Errorf("Channel %s is not of type *Channel, instead %T", channelID, tmp)
	}

	return channel, nil
}

func (h *Hub) printStats() {
	totalClients := 0

	log := config.Logger().With(
		"channels", h.channels.Count(),
	)
	for channelID, val := range h.channels.Items() {
		if channel, ok := val.(*Channel); ok {
			log = log.With(channelID, channel.clients.Count())
			totalClients += channel.clients.Count()
		}
	}

	log.Infow("Server Stats", "clients", totalClients)
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
	config.Logger().Infow(
		"Registered client to channel",
		"client_id", client.ID,
		"channel_id", client.channelID,
	)
}

// UnregisterClient ...
func (h *Hub) UnregisterClient(client *Client) {
	if client.channelID == uuid.Nil {
		config.Logger().DPanicw("Cant unregister client from nil channel id", "client_id", client.ID)
		return
	}

	channel, err := h.GetChannel(client.channelID)
	if err != nil {
		config.Logger().Warnw("Could not find channel to unregister client from", "client_id", client.ID, "channel_id", client.channelID, "reason", err.Error())
		return
	}

	// Remove client from channel & close the connection
	channel.clients.Remove(client.ID.String())

	// optionally, remove channel from hub if it has no clients left
	if channel.clients.Count() == 0 {
		h.channels.Remove(client.channelID.String())
		config.Logger().Debugw("Removing empty channel", "channel_id", client.channelID)
	}
}

func (h *Hub) SetGameIdForChannel(channelID uuid.UUID, gameID int) error {
	channel, err := h.GetChannel(channelID)
	if err != nil {
		return fmt.Errorf("Channel %s not found to set game id", channelID)
	}

	if channel.gameID != 0 && channel.gameID != gameID {
		return fmt.Errorf("Channel %s is already set to game %d, can not overwrite it to %d", channelID, channel.gameID, gameID)
	}

	channel.gameID = gameID
	return nil
}

func (h *Hub) ApplyProgressToChannel(payload *json.RawMessage, channelID uuid.UUID) (*json.RawMessage, error) {
	currentMsg, err := actions.NewApplyProgressMessageFrom(payload)

	if err != nil {
		return nil, fmt.Errorf("Could not parse apply message from payload to channel %s: %s", channelID, err)
	}

	tmp, ok := h.channels.Get(channelID.String())
	if !ok {
		config.Logger().Errorw("Channel not found to set apply progress to. Using request payload", "channel_id", channelID)
		return payload, nil
	}

	channel, ok := tmp.(*Channel)
	if !ok {
		config.Logger().Warnw("Couldnt cast to channel type to apply progress to. Using request payload", "channel_id", channelID)
		return payload, nil
	}

	latestMsg := actions.GetLatestProgress(channel.Progress, currentMsg)
	if currentMsg == latestMsg {
		config.Logger().Debug("New message was more recent than saved prgoress on channel")
	} else if latestMsg == channel.Progress {
		config.Logger().Debug("Current channel progress was more up to date")
	}

	channel.Progress = latestMsg

	data, err := json.Marshal(&channel.Progress)
	if err != nil {
		config.Logger().Errorw("Error marshalling apply progress message", "channel_id", channelID, "error", err.Error())
		return payload, nil
	}

	raw := json.RawMessage(data)
	return &raw, nil
}

// Broadcast ...
func (h *Hub) Broadcast(message *HubMessage) {
	from := message.client.ID.String()
	channelID := message.client.channelID.String()

	tmp, ok := h.channels.Get(channelID)
	if !ok {
		config.Logger().Debugw("Could not find channel to broadcast message to", "channel_id", channelID)
		return
	}

	channel, ok := tmp.(*Channel)
	if !ok {
		config.Logger().Debugw("Channel was not of type *Channel -- actually was %T", tmp)
		return
	}

	config.Logger().Debugw("Broadcasting message", "type", message.action)
	for to, _client := range channel.clients.Items() {
		// send to everyone else in the channel
		if to == from && !message.sendAll {
			continue
		}

		client, ok := _client.(*Client)
		if !ok {
			config.Logger().Warnf("Unknown client to broadcast to, skipping", "type", fmt.Sprintf("%T", _client), "client_id", to)
			continue
		}

		select {
		case client.send <- message.data:
		default:
			h.UnregisterClient(client)
		}

		config.Logger().Infow("Sent message to client", "client_id", to, "type", message.action)
	}
}
