package users

import (
	"github.com/graphql-go/graphql"
)

type UsersMutation struct {
	RegisterAccount *graphql.Field
	UpdateAccount   *graphql.Field
}

func InitUserMutation(r *UsersResolver) *UsersMutation {
	return &UsersMutation{
		RegisterAccount: &graphql.Field{
			Type:        userType,
			Description: "Register a new account",
			Resolve:     r.RegisterAccount,
		},
		UpdateAccount: &graphql.Field{
			Type:        userType,
			Description: "Update an account",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"imageURL": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"currentCharacterID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.UpdateAccount,
		},
	}
}
