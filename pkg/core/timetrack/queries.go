package timetrack

import (
	"github.com/graphql-go/graphql"
)

var TimeTrackingQuery = graphql.Field{
	Type:        timeTrackingType,
	Description: "Get a time tracking by ID",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: getTimeTrackingByID,
}

var AllTimeTrackingsQuery = graphql.Field{
	Type:        graphql.NewList(timeTrackingType),
	Description: "Get all time trackings",
	Resolve:     getAllTimeTrackings,
}

var TimeTrackingsByCharacterIDQuery = graphql.Field{
	Type:        graphql.NewList(timeTrackingType),
	Description: "Get all time trackings for a character by character ID",
	Args: graphql.FieldConfigArgument{
		"characterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: getTimeTrackingsByCharacterID,
}

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"timeTracking":               &TimeTrackingQuery,
		"timeTrackings":              &AllTimeTrackingsQuery,
		"timeTrackingsByCharacterID": &TimeTrackingsByCharacterIDQuery,
	},
})
