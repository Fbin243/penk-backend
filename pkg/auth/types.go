package auth

type ContextKey int

const (
	ProfileContextKey ContextKey = iota
	UserKey
)

type Profile struct {
	UID           string
	Email         string
	EmailVerified bool
	Name          string
}
