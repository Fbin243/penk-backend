package characters

import (
	"github.com/graphql-go/graphql"
)

type CharactersQuery struct {
	Character        *graphql.Field
	Characters       *graphql.Field
	UserCharacters   *graphql.Field
	CurrentCharacter *graphql.Field
}

func InitCharacterQuery(r *CharactersResolver) *CharactersQuery {
	return &CharactersQuery{
		Character: &graphql.Field{
			Type:        characterType,
			Description: "Get a character by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.GetCharacterByID,
		},
		Characters: &graphql.Field{
			Type:        graphql.NewList(characterType),
			Description: "Get all characters",
			Resolve:     r.GetAllCharacters,
		},
		UserCharacters: &graphql.Field{
			Type:        graphql.NewList(characterType),
			Description: "Get all characters of a user",
			Resolve:     r.GetCharactersByUserID,
		},
	}
}
