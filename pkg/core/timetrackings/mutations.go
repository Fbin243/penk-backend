package timetrackings

import (
	"github.com/graphql-go/graphql"
)

type TimeTrackingsMutation struct {
	CreateTimeTracking *graphql.Field
	UpdateTimeTracking *graphql.Field
}

func InitTimeTrackingsMutation(r *TimeTrackingsResolver) *TimeTrackingsMutation {
	return &TimeTrackingsMutation{
		CreateTimeTracking: &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Create a time tracking",
			Args: graphql.FieldConfigArgument{
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"customMetricID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.CreateTimeTracking,
		},
		UpdateTimeTracking: &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Update a time tracking",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.UpdateTimeTracking,
		},
	}
}
