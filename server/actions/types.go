package actions

import "strings"

type Action int

const (
	Register Action = iota
	SetGameId
	ApplyProgress
	Guess
	SelectCell
	OpenVeil
	CloseVeil
	Navigate
	SwitchDirection
	Clear
	ToggleHelpMenu
	EnterRebusMode
	ExitRebusMode

	Unknown
)

// IgnoredActions are actions we don't want to rebroadcast
var IgnoredActions = []Action{
	SelectCell,
	Navigate,
	SwitchDirection,
	ToggleHelpMenu,
	EnterRebusMode,
	ExitRebusMode,
}

var actionStrings = []string{
	"__CROSSWORD_PARTY_REGISTER",
	"__CROSSWORD_PARTY_SET_GAME_ID",
	"APPLY_PROGRESS",
	"GUESS",
	"SELECT_CELL",
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

func (a Action) String() string {
	return actionStrings[a]
}

func ActionFromString(str string) Action {
	str = strings.ToUpper(str)
	for index, actionStr := range actionStrings {
		if str == actionStr {
			return Action(index)
		}
	}
	return Unknown
}

func IsIgnoredAction(actionStr string) bool {
	a := ActionFromString(actionStr)

	for _, action := range IgnoredActions {
		if a == action {
			return true
		}
	}

	return false
}
