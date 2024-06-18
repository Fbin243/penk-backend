package characters

import (
	"github.com/graphql-go/graphql"
)

var CharacterQuery = graphql.Field{
	Type:        characterType,
	Description: "Get a character",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: getCharacterByID,
}

var CharactersQuery = graphql.Field{
	Type:        graphql.NewList(characterType),
	Description: "Get all characters",
	Resolve:     getAllCharacters,
}
