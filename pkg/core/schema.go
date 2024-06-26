package core

import (
	"tenkhours/pkg/core/characters"
	"tenkhours/pkg/core/timetrackings"
	"tenkhours/pkg/core/users"

	"github.com/graphql-go/graphql"
)

var (
	characterResolver = characters.NewCharactersResolver()
	characterQuery    = characters.InitCharacterQuery(characterResolver)
	characterMutation = characters.InitCharacterMutation(characterResolver)
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"user":           &users.User,
		"userCharacters": characterQuery.UserCharacters,
		"characters":     characterQuery.Characters,
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"registerAccount":    &users.RegisterAccount,
		"createCharacter":    characterMutation.CreateCharacter,
		"updateCharacter":    characterMutation.UpdateCharacter,
		"deleteCharacter":    characterMutation.DeleteCharacter,
		"resetCharacter":     characterMutation.ResetCharacter,
		"createCustomMetric": characterMutation.CreateCustomMetric,
		"updateCustomMetric": characterMutation.UpdateCustomMetric,
		"deleteCustomMetric": characterMutation.DeleteCustomMetric,
		"resetCustomMetric":  characterMutation.ResetCustomMetric,
		"createTimeTracking": &timetrackings.CreateTimeTrackingMutation,
		"updateTimeTracking": &timetrackings.UpdateTimeTrackingMutation,
	},
})

var CoreSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
