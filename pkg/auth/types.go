package auth

type ContextKey int

const (
	ProfileKey ContextKey = iota
	UserKey
	RegisterKey
	PostDataKey
)

type Profile struct {
	UID           string
	Email         string
	EmailVerified bool
	Name          string
}
