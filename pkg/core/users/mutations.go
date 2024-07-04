package users

import (
	"github.com/graphql-go/graphql"
)

type UsersMutation struct {
	UpdateAccount *graphql.Field
}

func InitUserMutation(r *UsersResolver) *UsersMutation {
	return &UsersMutation{
		UpdateAccount: &graphql.Field{
			Type:        userType,
			Description: "Update an account",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(userInputType),
				},
			},
			Resolve: r.UpdateAccount,
		},
	}
}
