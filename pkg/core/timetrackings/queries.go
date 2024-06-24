package timetrackings

import (
	"github.com/graphql-go/graphql"
)

var TimeTracking = graphql.Field{
	Type:        timeTrackingType,
	Description: "Get a time tracking by ID",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: getTimeTrackingByID,
}

var AllTimeTrackings = graphql.Field{
	Type:        graphql.NewList(timeTrackingType),
	Description: "Get all time trackings",
	Resolve:     getAllTimeTrackings,
}

var TimeTrackingsByCharacterID = graphql.Field{
	Type:        graphql.NewList(timeTrackingType),
	Description: "Get all time trackings for a character by character ID",
	Args: graphql.FieldConfigArgument{
		"characterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: getTimeTrackingsByCharacterID,
}
