package characters

import (
	"github.com/graphql-go/graphql"
)

type CharactersQuery struct {
	CharacterSnapshots *graphql.Field
	UserSnapshots      *graphql.Field
}

func InitCharacterQuery(r *CharactersResolver) *CharactersQuery {
	return &CharactersQuery{
		CharacterSnapshots: &graphql.Field{
			Type:        graphql.NewList(snapshotType),
			Description: "Get all snapshots of a character",
			Args: graphql.FieldConfigArgument{
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: r.GetSnapshotsByCharacterID,
		},
		UserSnapshots: &graphql.Field{
			Type:        graphql.NewList(snapshotType),
			Description: "Get all snapshots of all user's character",
			Resolve:     r.GetSnapshotsByUserID,
		},
	}
}
