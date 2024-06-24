package timetrackings

import (
	"github.com/graphql-go/graphql"
)

var CreateTimeTrackingMutation = graphql.Field{
	Type:        timeTrackingType,
	Description: "Create a time tracking",
	Args: graphql.FieldConfigArgument{
		"characterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"customMetricID": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: createTimeTracking,
}

var UpdateTimeTrackingMutation = graphql.Field{
	Type:        timeTrackingType,
	Description: "Update a time tracking",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: updateTimeTracking,
}
