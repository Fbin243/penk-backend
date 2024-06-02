package auth

import "fmt"

var (
	ErrorProfileNotFound    = fmt.Errorf("profile not found")
	ErrorCannotParseProfile = fmt.Errorf("cannot parse profile")
)
