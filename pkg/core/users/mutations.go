package users

import (
	"github.com/graphql-go/graphql"
)

type UserMutation struct {
	RegisterAccount *graphql.Field
}

func InitUserMutation(r *UsersResolver) *UserMutation {
	return &UserMutation{
		RegisterAccount: &graphql.Field{
			Type:        userType,
			Description: "Register a new account",
			Resolve:     r.RegisterAccount,
		},
	}
}
