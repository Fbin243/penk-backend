package gql

import (
	"tenkhours/pkg/core"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/db/timetrackingsdb"

	"github.com/graphql-go/graphql"
)

func InitSchema() graphql.Schema {
	var (
		db                = db.GetDBManager().DB
		charactersRepo    = coredb.NewCharactersRepo(db)
		timeTrackingsRepo = timetrackingsdb.NewTimeTrackingsRepo(db)

		timeTrackingsResolver = core.NewTimeTrackingsResolver(timeTrackingsRepo, charactersRepo)
		timeTrackingsQuery    = InitTimeTrackingsQuery(timeTrackingsResolver)
		timeTrackingsMutation = InitTimeTrackingsMutation(timeTrackingsResolver)
	)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"currentTimeTracking": timeTrackingsQuery.CurrentTimeTracking,
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createTimeTracking": timeTrackingsMutation.CreateTimeTracking,
			"updateTimeTracking": timeTrackingsMutation.UpdateTimeTracking,
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	return schema
}
