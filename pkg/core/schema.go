package core

import (
	"tenkhours/pkg/core/characters"
	"tenkhours/pkg/core/timetrackings"
	"tenkhours/pkg/core/users"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
)

func InitSchema() graphql.Schema {
	var (
		db                = db.GetDBManager().DB
		usersRepo         = coredb.NewUsersRepo(db)
		charactersRepo    = coredb.NewCharactersRepo(db)
		timeTrackingsRepo = coredb.NewTimeTrackingsRepo(db)

		usersResolver = users.NewUsersResolver(usersRepo)
		usersQuery    = users.InitUserQuery(usersResolver)
		usersMutation = users.InitUserMutation(usersResolver)

		charactersResolver = characters.NewCharactersResolver(charactersRepo, usersRepo)
		charactersQuery    = characters.InitCharacterQuery(charactersResolver)
		charactersMutation = characters.InitCharacterMutation(charactersResolver)

		timeTrackingsResolver = timetrackings.NewTimeTrackingsResolver(timeTrackingsRepo, charactersRepo)
		timeTrackingsQuery    = timetrackings.InitTimeTrackingsQuery(timeTrackingsResolver)
		timeTrackingsMutation = timetrackings.InitTimeTrackingsMutation(timeTrackingsResolver)
	)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user":                usersQuery.User,
			"userCharacters":      charactersQuery.UserCharacters,
			"characters":          charactersQuery.Characters,
			"currentTimeTracking": timeTrackingsQuery.CurrentTimeTracking,
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"updateAccount": usersMutation.UpdateAccount,

			"createCharacter": charactersMutation.CreateCharacter,
			"updateCharacter": charactersMutation.UpdateCharacter,
			"deleteCharacter": charactersMutation.DeleteCharacter,
			"resetCharacter":  charactersMutation.ResetCharacter,

			"createCustomMetric": charactersMutation.CreateCustomMetric,
			"updateCustomMetric": charactersMutation.UpdateCustomMetric,
			"deleteCustomMetric": charactersMutation.DeleteCustomMetric,
			"resetCustomMetric":  charactersMutation.ResetCustomMetric,

			"createMetricProperty": charactersMutation.CreateMetricProperty,
			"updateMetricProperty": charactersMutation.UpdateMetricProperty,
			"deleteMetricProperty": charactersMutation.DeleteMetricProperty,

			"createTimeTracking": timeTrackingsMutation.CreateTimeTracking,
			"updateTimeTracking": timeTrackingsMutation.UpdateTimeTracking,
		},
	})

	CoreSchema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	return CoreSchema
}
