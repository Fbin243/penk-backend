package gql

import (
	"tenkhours/pkg/core"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
)

func InitSchema() graphql.Schema {
	var (
		db             = db.GetDBManager().DB
		usersRepo      = coredb.NewUsersRepo(db)
		charactersRepo = coredb.NewCharactersRepo(db)

		usersResolver = core.NewUsersHandler(usersRepo)
		usersQuery    = InitUserQuery(usersResolver)
		usersMutation = InitUserMutation(usersResolver)

		charactersResolver = core.NewCharactersHandler(charactersRepo, usersRepo)
		charactersQuery    = InitCharacterQuery(charactersResolver)
		charactersMutation = InitCharacterMutation(charactersResolver)
	)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user":           usersQuery.User,
			"userCharacters": charactersQuery.UserCharacters,
			"characters":     charactersQuery.Characters,
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
		},
	})

	CoreSchema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	return CoreSchema
}
