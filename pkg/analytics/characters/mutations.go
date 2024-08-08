package characters

import "github.com/graphql-go/graphql"

type CharactersMutation struct {
	CreateSnapshot *graphql.Field
}

func InitCharactersMutation(r *CharactersResolver) *CharactersMutation {
	return &CharactersMutation{
		CreateSnapshot: &graphql.Field{
			Type: snapshotType,
			Args: graphql.FieldConfigArgument{
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: r.CreateNewSnapshot,
		},
	}
}
