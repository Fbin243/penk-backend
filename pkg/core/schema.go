package core

import (
	"tenkhours/pkg/core/characters"
	"tenkhours/pkg/core/timetrackings"
	"tenkhours/pkg/core/users"

	"github.com/graphql-go/graphql"
)

func InitSchema() graphql.Schema {
	var (
		usersResolver     = users.NewUsersResolver()
		userQuery         = users.InitUserQuery(usersResolver)
		userMutation      = users.InitUserMutation(usersResolver)
		characterResolver = characters.NewCharactersResolver()
		characterQuery    = characters.InitCharacterQuery(characterResolver)
		characterMutation = characters.InitCharacterMutation(characterResolver)
	)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user":           userQuery.User,
			"userCharacters": characterQuery.UserCharacters,
			"characters":     characterQuery.Characters,
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"registerAccount":    userMutation.RegisterAccount,
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

	CoreSchema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	return CoreSchema
}
