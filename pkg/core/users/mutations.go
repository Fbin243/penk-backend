package users

import (
	"github.com/graphql-go/graphql"
)

type UsersMutation struct {
	RegisterAccount *graphql.Field
}

func InitUserMutation(r *UsersResolver) *UsersMutation {
	return &UsersMutation{
		RegisterAccount: &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Register a new account",
			Resolve:     r.RegisterAccount,
		},
	}
}
