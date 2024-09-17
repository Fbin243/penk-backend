package gql

import (
	"tenkhours/pkg/analytics"
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

		charactersResolver = analytics.NewCharactersResolver(snapshotsRepo, charactersRepo, usersRepo)
		charactersQuery    = InitCharacterQuery(charactersResolver)
		charactersMutation = InitCharactersMutation(charactersResolver)
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

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	return schema
}
