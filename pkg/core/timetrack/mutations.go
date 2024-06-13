package timetrack

import (
	"github.com/graphql-go/graphql"
)

// var newTimeTrackingInput = graphql.NewInputObject(graphql.InputObjectConfig{
// 	Name: "TimeTrackingInput",
// 	Fields: graphql.InputObjectConfigFieldMap{
// 		"characterID": &graphql.InputObjectFieldConfig{
// 			Type: graphql.String,
// 		},
// 		"startTime": &graphql.InputObjectFieldConfig{
// 			Type: graphql.Int,
// 		},
// 		"endTime": &graphql.InputObjectFieldConfig{
// 			Type: graphql.Int,
// 		},
// 	},
// })
// resolver functions are directly accessing the individual arguments from the params.Args map

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createTimeTracking": &graphql.Field{
			Type:        timeTrackingType,
			Description: "Create a time tracking",
			Args: graphql.FieldConfigArgument{
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"startTime": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"endTime": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: createTimeTracking,
		},

		"updateTimeTracking": &graphql.Field{
			Type:        timeTrackingType,
			Description: "Update a time tracking",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"startTime": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"endTime": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: updateTimeTracking,
		},

		"deleteTimeTracking": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete a time tracking",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: deleteTimeTracking,
		},
	},
})