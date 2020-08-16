package actions

type Action string

const (
	Register        Action = "__CROSSWORD_PARTY_REGISTER"
	SetGameID              = "__CROSSWORD_PARTY_SET_GAME_ID"
	ApplyProgress          = "APPLY_PROGRESS"
	Guess                  = "GUESS"
	SelectCell             = "SELECT_CELL"
	SelectClue             = "SELECT_CLUE"
	OpenVeil               = "OPEN_VEIL"
	CloseVeil              = "CLOSE_VEIL"
	Navigate               = "NAVIGATE"
	SwitchDirection        = "SWITCH_DIRECTION"
	Clear                  = "CLEAR"
	ToggleHelpMenu         = "TOGGLE_HELP_MENU"
	EnterRebusMode         = "ENTER_REBUS_MODE"
	ExitRebusMode          = "EXIT_REBUS_MODE"
	Unknown                = "UNKNOWN"
)

// Actions is a mapping of which actions to ignore
var Actions = map[Action]bool{
	// Actions to rebroadcast explicitly
	ApplyProgress: false,
	Guess:         false,
	OpenVeil:      false,
	CloseVeil:     false,
	Clear:         false,
	Unknown:       false,

	// Actioans to ignore during rebroadcasting
	Register:        true,
	SetGameID:       true,
	SelectCell:      true,
	SelectClue:      true,
	Navigate:        true,
	SwitchDirection: true,
	ToggleHelpMenu:  true,
	EnterRebusMode:  true,
	ExitRebusMode:   true,
}

var actionStrings = []string{
	"__CROSSWORD_PARTY_REGISTER",
	"__CROSSWORD_PARTY_SET_GAME_ID",
	"APPLY_PROGRESS",
	"GUESS",
	"SELECT_CELL",
	"SELECT_CLUE",
	"OPEN_VEIL",
	"CLOSE_VEIL",
	"NAVIGATE",
	"SWITCH_DIRECTION",
	"CLEAR",
	"TOGGLE_HELP_MENU",
	"ENTER_REBUS_MODE",
	"EXIT_REBUS_MODE",
	"UNKNOWN",
}

func IsIgnoredAction(action Action) bool {
	return Actions[action]
}
