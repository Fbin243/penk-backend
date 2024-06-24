package core

import (
	"tenkhours/pkg/core/characters"
	"tenkhours/pkg/core/users"
	"tenkhours/pkg/core/timetrack"

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
		"registerAccount":    &users.RegisterAccount,
		"createCharacter":    &characters.CreateCharacter,
		"updateCharacter":    &characters.UpdateCharacter,
		"deleteCharacter":    &characters.DeleteCharacter,
		"resetCharacter":     &characters.ResetCharacter,
		"createCustomMetric": &characters.CreateCustomMetric,
		"updateCustomMetric": &characters.UpdateCustomMetric,
		"deleteCustomMetric": &characters.DeleteCustomMetric,
		"resetCustomMetric":  &characters.ResetCustomMetric,
	},
})

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"timeTracking":               &timetrack.TimeTrackingQuery,
		"timeTrackings":              &timetrack.AllTimeTrackingsQuery,
		"timeTrackingsByCharacterID": &timetrack.TimeTrackingsByCharacterIDQuery,
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createTimeTracking": &timetrack.CreateTimeTrackingMutation,
		"updateTimeTracking": &timetrack.UpdateTimeTrackingMutation,
		"deleteTimeTracking": &timetrack.DeleteTimeTrackingMutation,
	},
})


var CoreSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
