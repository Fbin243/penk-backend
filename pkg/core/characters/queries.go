package characters

import (
	"github.com/graphql-go/graphql"
)

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"character": &graphql.Field{
			Type:        characterType,
			Description: "Get a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: getCharacterByID,
		},
		"characters": &graphql.Field{
			Type:        graphql.NewList(characterType),
			Description: "Get all characters",
			Resolve:     getAllCharacters,
		},
	},
})
