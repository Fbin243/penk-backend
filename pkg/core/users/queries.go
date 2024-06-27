package users

import (
	"github.com/graphql-go/graphql"
)

type UserQuery struct {
	User *graphql.Field
}

func InitUserQuery(r *UsersResolver) *UserQuery {
	return &UserQuery{
		User: &graphql.Field{
			Type:        userType,
			Description: "Get a user by email",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.GetUserByEmail,
		},
	}
}
