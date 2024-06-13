package characters

import (
	"github.com/graphql-go/graphql"
)

var newCustomMetricInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CusTomMetricInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"characterID": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"type": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"value": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createCharacter": &graphql.Field{
			Type:        characterType,
			Description: "Create a character",
			Args: graphql.FieldConfigArgument{
				"userID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"tags": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
				"totalFocusTime": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"customMetrics": &graphql.ArgumentConfig{
					Type: graphql.NewList(newCustomMetricInput),
				},
			},
			Resolve: createCharacter,
		},

		"updateCharacter": &graphql.Field{
			Type:        characterType,
			Description: "Update a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"userID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"tags": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
				"totalFocusTime": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"customMetrics": &graphql.ArgumentConfig{
					Type: graphql.NewList(newCustomMetricInput),
				},
			},
			Resolve: updateCharacter,
		},

		"deleteCharacter": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: deleteCharacter,
		},
	},
})
