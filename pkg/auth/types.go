package auth

type ContextKey int

const (
	ProfileKey ContextKey = iota
)

type FirebaseProfile struct {
	UID     string
	Email   string
	Name    string
	Picture string
}

type AuthOption struct {
	Authorized bool
}
