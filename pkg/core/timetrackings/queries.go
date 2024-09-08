package timetrackings

import (
	"github.com/graphql-go/graphql"
)

type TimeTrackingsQuery struct {
	CurrentTimeTracking *graphql.Field
}

func InitTimeTrackingsQuery(r *TimeTrackingsResolver) *TimeTrackingsQuery {
	return &TimeTrackingsQuery{
		CurrentTimeTracking: &graphql.Field{
			Type:        timeTrackingType,
			Description: "Get current time tracking (can be null)",
			Args: graphql.FieldConfigArgument{
				"characterID": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.ID),
					Description: "ID character to track time for",
				},
			},
			Resolve: r.GetCurrentTimeTracking,
		},
	}
}
