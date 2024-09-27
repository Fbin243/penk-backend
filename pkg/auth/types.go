package auth

type ContextKey int

const (
	ProfileKey ContextKey = iota
	PostDataKey
)

type Profile struct {
	UID           string
	Email         string
	EmailVerified bool
	Name          string
}
