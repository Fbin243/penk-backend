package users

import (
	"github.com/graphql-go/graphql"
)

var GetUserById = graphql.Field{
	Type:        userType,
	Description: "Get a user by ID",
}
