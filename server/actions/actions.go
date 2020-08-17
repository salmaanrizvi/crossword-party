package actions

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//Message for generic redux message
type Message struct {
	Type          Action    `json:"type"`
	From          uuid.UUID `json:"from"`
	Channel       uuid.UUID `json:"channel"`
	Timestamp     time.Time `json:"timestamp"`
	ClientVersion string    `json:"clientVersion"`
	GameID        int       `json:"gameId"`
}

//NewMessageFrom parses websocket message into a Message struct
func NewMessageFrom(data []byte) (msg *Message, err error) {
	msg = &Message{}
	err = json.Unmarshal(data, msg)
	return msg, err
}
