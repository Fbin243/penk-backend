package gql

import (
	"tenkhours/pkg/core"

	"github.com/graphql-go/graphql"
)

type ProfilesMutation struct {
	UpdateAccount *graphql.Field
}

func InitProfileMutation(r *core.ProfilesHandler) *ProfilesMutation {
	return &ProfilesMutation{
		UpdateAccount: &graphql.Field{
			Type:        profileType,
			Description: "Update an account",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(profileInputType),
				},
			},
			Resolve: r.UpdateAccount,
		},
	}
}
