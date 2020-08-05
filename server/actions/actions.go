package actions

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//Message for generic redux message
type Message struct {
	Type          string    `json:"type"`
	From          uuid.UUID `json:"from"`
	Channel       uuid.UUID `json:"channel"`
	Timestamp     time.Time `json:"timestamp"`
	ClientVersion string    `json:"clientVersion"`
	GameId        int       `json:"gameId"`
}

//NewMessageFrom parses websocket message into a Message struct
func NewMessageFrom(data []byte) (msg *Message, err error) {
	msg = &Message{}
	err = json.Unmarshal(data, msg)
	return msg, err
}

type ApplyProgressMessage struct {
	Message
	Selection *Selection            `json:"selection"`
	Payload   *ApplyProgressPayload `json:"payload"`
}

func NewApplyProgessMessageFrom(data []byte) (apMsg *ApplyProgressMessage, err error) {
	apMsg = &ApplyProgressMessage{}
	err = json.Unmarshal(data, apMsg)
	return apMsg, err
}

type ApplyProgressPayload struct {
	Cells     []*Cell    `json:"cells"`
	Selection *Selection `json:"selection"`
}

type Selection struct {
	Cell         int     `json:"cell"`
	Clue         int     `json:"clue"`
	ClueList     int     `json:"clueList"`
	CellClues    []uint8 `json:"cellClues"`
	ClueCells    []uint8 `json:"clueCells"`
	RelatedCells []uint8 `json:"relatedCells"`
	RelatedClues []uint8 `json:"relatedClues"`
}

type Cell struct {
	Answer    string  `json:"answer"`
	Clues     []uint8 `json:"clues"`
	Confirmed bool    `json:"confirmed"`
	Guess     string  `json:"guess"`
	Index     int     `json:"index"`
	Label     string  `json:"label"`
	Penciled  bool    `json:"penciled"`
	Timestamp int     `json:"timestamp"`
	CellType  int     `json:"type"`
}

type SelectCellMessage struct {
	Message
	Payload *SelectCellPayload `json:"payload"`
}

type SelectCellPayload struct {
	Index         int  `json:"index"`
	isMiddleClick bool `json:"isMiddleClick"`
}

type GuessMessage struct {
	Message
	Payload *GuessPayload `json:"payload"`
}

type GuessPayload struct {
	BlankDelta       int    `json:"blankDelta"`
	IncorrectDelta   int    `json:"incorrectDelta"`
	Index            int    `json:"index"`
	InPencilMode     bool   `json:"inPencilMode"`
	AutocheckEnabled bool   `json:"autocheckEnabled"`
	Value            string `json:"value"`
	FromRebus        bool   `json:"fromRebus"`
	Now              int    `json:"now"`
}

func NewGuessMessageFrom(data []byte) (guess *GuessMessage, err error) {
	guess = &GuessMessage{}
	err = json.Unmarshal(data, guess)
	return guess, err
}
