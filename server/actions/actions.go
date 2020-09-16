package actions

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/salmaanrizvi/crossword-party/server/config"
)

//Message for generic redux message
type Message struct {
	Type          Action          `json:"type"`
	From          uuid.UUID       `json:"from"`
	Channel       uuid.UUID       `json:"channel"`
	Timestamp     time.Time       `json:"timestamp"`
	ClientVersion string          `json:"clientVersion"`
	GameID        int             `json:"gameId"`
	IsFromServer  bool            `json:"isFromServer"`
	Payload       json.RawMessage `json:"payload"`
}

//NewMessageFrom parses websocket message into a Message struct
func NewMessageFrom(data []byte) (msg *Message, err error) {
	msg = &Message{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	msg.IsFromServer = true
	return msg, nil
}

func (msg *Message) RawMessage() []byte {
	start := time.Now()
	defer func() {
		config.Logger().Debugf("RawMessage parsing took %v", time.Now().Sub(start))
	}()

	data, err := json.Marshal(msg)
	if err != nil {
		return []byte("{}")
	}
	return data
}

// ApplyProgressMessage represents the "APPLY_PROGRESS" redux action
// sent by clients to set the game board with filled in cells & timer
type ApplyProgressMessage struct {
	Cells  []Cell  `json:"cells"`
	Status *Status `json:"status"`
	Timer  *Timer  `json:"timer"`
}

// Cell represents a single cell on a game board
type Cell struct {
	Index     int    `json:"index"`
	Clues     []int  `json:"clues"`
	Answer    string `json:"answer"`
	Checked   bool   `json:"checked"`
	Confirmed bool   `json:"confirmed"`
	Guess     string `json:"guess"`
	Label     string `json:"label"`
	Modified  bool   `json:"modified"`
	Penciled  bool   `json:"penciled"`
	Revealed  bool   `json:"revealed"`
}

// Status represents the games active status
type Status struct {
	AutoCheckEnabled bool            `json:"autocheckEnabled"`
	BlankCells       int             `json:"blankCells"`
	Firsts           *Firsts         `json:"firsts"`
	IncorrectCells   int             `json:"incorrectCells"`
	IsFilled         bool            `json:"isFilled"`
	IsSolved         bool            `json:"isSolved"`
	LastCommitID     string          `json:"lastCommitID"`
	CurrentProgress  json.RawMessage `json:"currentProgress"`
}

// Firsts is... unclear if this is useful
type Firsts struct {
	Opened int64 `json:"opened"`
}

// Timer is the state of the game's timer
type Timer struct {
	ResetSinceLastCommit bool `json:"resetSinceLastCommit"`
	SessionElapsedTime   int  `json:"sessionElapsedTime"`
	TotalElapsedTime     int  `json:"totalElapsedTime"`
}

// NewApplyProgressMessageFrom returns a new ApplyProgressMessage from raw data received over the websocket
func NewApplyProgressMessageFrom(data json.RawMessage) (apMsg *ApplyProgressMessage, err error) {
	apMsg = &ApplyProgressMessage{}
	err = json.Unmarshal(data, apMsg)
	return apMsg, err
}

func (ap *ApplyProgressMessage) RawMessage() []byte {
	start := time.Now()
	data, err := json.Marshal(ap)
	defer func() {
		config.Logger().Debugf("Apply Progress RawMessage parsing took %v", time.Now().Sub(start))
		config.Logger().Debugf("RAW AP: %s", data)
	}()
	if err != nil {
		return []byte("{}")
	}
	return data
}

// GetLatestProgress compares progress messages and returns the one which
// has more of the puzzle completed
func GetLatestProgress(a, b *ApplyProgressMessage) *ApplyProgressMessage {
	if a == nil && b == nil {
		return nil
	} else if a != nil && b == nil {
		return a
	} else if a == nil && b != nil {
		return b
	}

	if a.Status.IsSolved && b.Status.IsSolved {
		return a
	} else if a.Status.IsSolved && !b.Status.IsSolved {
		return a
	} else if !a.Status.IsSolved && b.Status.IsSolved {
		return b
	}

	if a.Status.IncorrectCells < b.Status.IncorrectCells {
		return a
	} else if b.Status.IncorrectCells < a.Status.IncorrectCells {
		return b
	}

	if a.Status.BlankCells < b.Status.BlankCells {
		return a
	} else if b.Status.BlankCells < a.Status.BlankCells {
		return b
	}

	return a
}
