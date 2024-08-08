package analytics

import (
	"tenkhours/pkg/analytics/characters"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/analyticsdb"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
)

func InitSchema() graphql.Schema {
	var (
		db             = db.GetDBManager().DB
		snapshotsRepo  = analyticsdb.NewSnapshotRepo(db)
		usersRepo      = coredb.NewUsersRepo(db)
		charactersRepo = coredb.NewCharactersRepo(db)

		charactersResolver = characters.NewCharactersResolver(snapshotsRepo, charactersRepo, usersRepo)
		charactersQuery    = characters.InitCharacterQuery(charactersResolver)
		charactersMutation = characters.InitCharactersMutation(charactersResolver)
	)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"characterSnapshots": charactersQuery.CharacterSnapshots,
			"userSnapshots":      charactersQuery.UserSnapshots,
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createSnapshot": charactersMutation.CreateSnapshot,
		},
	})

	AnalyticsSchema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	return AnalyticsSchema
}
