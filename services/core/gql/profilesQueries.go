package gql

import (
	"tenkhours/pkg/core"

	"github.com/graphql-go/graphql"
)

type ProfilesQuery struct {
	Profile *graphql.Field
}

func InitProfileQuery(r *core.ProfilesHandler) *ProfilesQuery {
	return &ProfilesQuery{
		Profile: &graphql.Field{
			Type:        graphql.NewNonNull(profileType),
			Description: "Get user's profile by token",
			Resolve:     r.GetProfileByToken,
		},
	}
}
