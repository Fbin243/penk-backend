package core

import (
<<<<<<< HEAD
	"tenkhours/pkg/core/characters"
=======
>>>>>>> dev
	"tenkhours/pkg/core/users"

	"github.com/graphql-go/graphql"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
<<<<<<< HEAD
		"user":       &users.User,
		"character":  &characters.CharacterQuery,
		"characters": &characters.CharactersQuery,
=======
		"user": &users.User,
>>>>>>> dev
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"registerAccount": &users.RegisterAccount,
<<<<<<< HEAD
		"createCharacter": &characters.CreateCharacter,
		"updateCharacter": &characters.UpdateCharacter,
		"deleteCharacter": &characters.DeleteCharacter,
=======
>>>>>>> dev
	},
})

var CoreSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
