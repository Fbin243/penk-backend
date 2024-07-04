package characters

import (
	"github.com/graphql-go/graphql"
)

type CharactersMutation struct {
	CreateCharacter *graphql.Field
	UpdateCharacter *graphql.Field
	DeleteCharacter *graphql.Field
	ResetCharacter  *graphql.Field

	CreateCustomMetric *graphql.Field
	UpdateCustomMetric *graphql.Field
	DeleteCustomMetric *graphql.Field
	ResetCustomMetric  *graphql.Field

	CreateMetricProperty *graphql.Field
	UpdateMetricProperty *graphql.Field
	DeleteMetricProperty *graphql.Field
}

func InitCharacterMutation(r *CharactersResolver) *CharactersMutation {
	return &CharactersMutation{
		CreateCharacter: &graphql.Field{
			Type:        characterType,
			Description: "Create a character",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"gender": &graphql.ArgumentConfig{
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
				"gender": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"tags": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
			},
			Resolve: r.UpdateCharacter,
		},
		DeleteCharacter: &graphql.Field{
			Type:        characterType,
			Description: "Delete a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.DeleteCharacter,
		},
		ResetCharacter: &graphql.Field{
			Type:        characterType,
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
			Description: "Create a custom metric",
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
					Type: metricStyleInput,
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
					Type: metricStyleInput,
				},
			},
			Resolve: r.UpdateCustomMetric,
		},
		ResetCustomMetric: &graphql.Field{
			Type:        customMetricType,
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
			Type:        customMetricType,
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
		CreateMetricProperty: &graphql.Field{
			Type:        metricPropertyType,
			Description: "Create a metric property",
			Args: graphql.FieldConfigArgument{
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"metricID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"type": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"value": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"unit": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.CreateMetricProperty,
		},
		UpdateMetricProperty: &graphql.Field{
			Type:        metricPropertyType,
			Description: "Update a metric property",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"metricID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"type": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"value": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"unit": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.UpdateMetricProperty,
		},
		DeleteMetricProperty: &graphql.Field{
			Type:        metricPropertyType,
			Description: "Delete a metric property",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"metricID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.DeleteMetricProperty,
		},
	}
}
