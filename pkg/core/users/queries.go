package users

import (
	"github.com/graphql-go/graphql"
)

var User = graphql.Field{
	Type:        userType,
	Description: "Get a user by email",
	Args: graphql.FieldConfigArgument{
		"email": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: getUserByEmail,
}
