package characters

import (
	"github.com/graphql-go/graphql"
)

type CharacterMutation struct {
	CreateCharacter *graphql.Field
	UpdateCharacter *graphql.Field
	DeleteCharacter *graphql.Field
	ResetCharacter  *graphql.Field

	CreateCustomMetric *graphql.Field
	UpdateCustomMetric *graphql.Field
	DeleteCustomMetric *graphql.Field
	ResetCustomMetric  *graphql.Field
}

func InitCharacterMutation(r *CharactersResolver) *CharacterMutation {
	return &CharacterMutation{
		CreateCharacter: &graphql.Field{
			Type:        characterType,
			Description: "Create a character",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"tags": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
			},
			Resolve: r.CreateCharacter,
		},
		UpdateCharacter: &graphql.Field{
			Type:        characterType,
			Description: "Update a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"tags": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
			},
			Resolve: r.UpdateCharacter,
		},
		DeleteCharacter: &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.DeleteCharacter,
		},
		ResetCharacter: &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Reset a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.ResetCharacter,
		},
		CreateCustomMetric: &graphql.Field{
			Type:        customMetricType,
			Description: "Create a Custom Metric",
			Args: graphql.FieldConfigArgument{
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"style": &graphql.ArgumentConfig{
					Type: metricStyleTypeInput,
				},
			},
			Resolve: r.CreateCustomMetric,
		},
		UpdateCustomMetric: &graphql.Field{
			Type:        customMetricType,
			Description: "Update a Custom Metric",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"style": &graphql.ArgumentConfig{
					Type: metricStyleTypeInput,
				},
				"properties": &graphql.ArgumentConfig{
					Type: graphql.NewList(metricPropertyTypeInput),
				},
			},
			Resolve: r.UpdateCustomMetric,
		},
		ResetCustomMetric: &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Reset a custom metric",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.ResetCustomMetric,
		},
		DeleteCustomMetric: &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete a custom metric",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.DeleteCustomMetric,
		},
	}
}
