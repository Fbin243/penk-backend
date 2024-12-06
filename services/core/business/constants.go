package business

type ContextKey int

const (
	CharacterKey ContextKey = iota
	FromCreateCharacter
	FromUpdateCharacter
)
