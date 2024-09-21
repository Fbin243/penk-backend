package gql

import (
	"tenkhours/pkg/core"

	"github.com/graphql-go/graphql"
)

type UsersQuery struct {
	User *graphql.Field
}

func InitUserQuery(r *core.UsersHandler) *UsersQuery {
	return &UsersQuery{
		User: &graphql.Field{
			Type:        graphql.NewNonNull(userType),
			Description: "Get a user by token",
			Resolve:     r.GetUserByToken,
		},
	}
}
