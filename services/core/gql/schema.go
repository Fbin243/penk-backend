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
		profilesRepo   = coredb.NewProfilesRepo(db)
		charactersRepo = coredb.NewCharactersRepo(db)

		profilesResolver = core.NewProfilesHandler(profilesRepo)
		profilesQuery    = InitProfileQuery(profilesResolver)
		profilesMutation = InitProfileMutation(profilesResolver)

		charactersResolver = core.NewCharactersHandler(charactersRepo, profilesRepo)
		charactersQuery    = InitCharacterQuery(charactersResolver)
		charactersMutation = InitCharacterMutation(charactersResolver)
	)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"profile":        profilesQuery.Profile,
			"userCharacters": charactersQuery.UserCharacters,
			"characters":     charactersQuery.Characters,
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"updateAccount": profilesMutation.UpdateAccount,

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
