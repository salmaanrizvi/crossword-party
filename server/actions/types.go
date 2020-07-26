package actions

type Action int

const (
	Register Action = iota
	ApplyProgress
	Guess
)

func (a Action) String() string {
	return []string{"__CROSSWORD_PARTY_REGISTER", "APPLY_PROGRESS", "GUESS"}[a]
}
