package gql

import (
	"tenkhours/pkg/core"

	"github.com/graphql-go/graphql"
)

type TimeTrackingsMutation struct {
	CreateTimeTracking *graphql.Field
	UpdateTimeTracking *graphql.Field
}

func InitTimeTrackingsMutation(r *core.TimeTrackingsHandler) *TimeTrackingsMutation {
	return &TimeTrackingsMutation{
		CreateTimeTracking: &graphql.Field{
			Type:        graphql.NewNonNull(timeTrackingType),
			Description: "Create a time tracking",
			Args: graphql.FieldConfigArgument{
				"characterID": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.ID),
					Description: "ID character to track time for",
				},
				"metricID": &graphql.ArgumentConfig{
					Type:        graphql.ID,
					Description: "ID custom metric to track time for",
				},
				"startTime": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.DateTime),
					Description: "The start time of new session in ISO 8601 format",
				},
			},
			Resolve: r.CreateTimeTracking,
		},
		UpdateTimeTracking: &graphql.Field{
			Type:        graphql.NewNonNull(timeTrackingType),
			Description: "Update a time tracking",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.ID),
					Description: "ID of time tracking to update",
				},
			},
			Resolve: r.UpdateTimeTracking,
		},
	}
}
