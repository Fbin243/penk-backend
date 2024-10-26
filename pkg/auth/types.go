package auth

type ContextKey int

const (
	FirebaseProfileKey ContextKey = iota
	ProfileKey
	PostDataKey
)

type FirebaseProfile struct {
	UID   string
	Email string
	Name  string
}

type AuthOption struct {
	Authorized bool
}
