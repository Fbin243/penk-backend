package core

import (
	"tenkhours/pkg/core/characters"
	"tenkhours/pkg/core/users"

	"github.com/graphql-go/graphql"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"user":       &users.User,
		"character":  &characters.CharacterQuery,
		"characters": &characters.CharactersQuery,
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"registerAccount": &users.RegisterAccount,
		"createCharacter": &characters.CreateCharacter,
		"updateCharacter": &characters.UpdateCharacter,
		"deleteCharacter": &characters.DeleteCharacter,
		"resetCharacter":  &characters.ResetCharacter,
	},
})

var CoreSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
