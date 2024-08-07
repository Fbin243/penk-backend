package common

type ContextKey string

const (
	TestingT         ContextKey = "TestingT"
	User             ContextKey = "User"
	AnotherCharacter ContextKey = "AnotherCharacter"
	CurrentCharacter ContextKey = "CurrentCharacter"
	Snapshot         ContextKey = "Snapshot"
)
