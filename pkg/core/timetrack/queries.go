package timetrack

import (
	"github.com/graphql-go/graphql"
)

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"timeTracking": &graphql.Field{
			Type:        timeTrackingType,
			Description: "Get a time tracking",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: getTimeTrackingByID,
		},

		"timeTrackings": &graphql.Field{
			Type:        graphql.NewList(timeTrackingType),
			Description: "Get all time trackings",
			Resolve:     getAllTimeTrackings,
		},

		"timeTrackingsByCharacterID": &graphql.Field{
			Type:        graphql.NewList(timeTrackingType),
			Description: "Get all time trackings for a character",
			Args: graphql.FieldConfigArgument{
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: getTimeTrackingsByCharacterID,
		},
	},
})