package users

import (
	"github.com/graphql-go/graphql"
)

var RegisterAccount = graphql.Field{
	Type:        userType,
	Description: "Register a new account",
	Resolve:     registerAccount,
}
