package types

import "encoding/json"

//Message for generic redux message
type Message struct {
	messageType string `json:"type"`
	payload     []byte `json:"payload"`
}

//NewMessageFrom parses websocket message into a Message struct
func NewMessageFrom(data []byte) (msg *Message, err error) {
	err = json.Unmarshal(data, msg)
	return msg, err
}

func (msg *Message) GetApplyProgress() (ap *ApplyProgress, err error) {
	err = json.Unmarshal(msg.payload, ap)
	return ap, err
}

// type Guess struct {

// }

type ApplyProgress struct {
	selection *Selection `json:"selection"`
	cells     []*Cell    `json:"cells"`
}

type Selection struct {
	cell         int   `json:"cell"`
	clue         int   `json:"clue"`
	clueList     int   `json:"clueList"`
	cellClues    []int `json:"cellClues"`
	clueCells    []int `json:"clueCells"`
	relatedCells []int `json:"relatedCells"`
	relatedClues []int `json:"relatedClues"`
}

type Cell struct {
	answer    string `json:"answer"`
	clues     []int  `json:"clues"`
	confirmed bool   `json:"confirmed"`
	guess     string `json:"guess"`
	index     int    `json:"index"`
	label     string `json:"label"`
	penciled  bool   `json:"penciled"`
	timestamp int    `json:"timestamp"`
	cellType  int    `json:"type"`
}
