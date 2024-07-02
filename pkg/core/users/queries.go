package users

import (
	"github.com/graphql-go/graphql"
)

type UsersQuery struct {
	User *graphql.Field
}

func InitUserQuery(r *UsersResolver) *UsersQuery {
	return &UsersQuery{
		User: &graphql.Field{
			Type:        userType,
			Description: "Get a user by token",
			Resolve:     r.GetUserByToken,
		},
	}
}
