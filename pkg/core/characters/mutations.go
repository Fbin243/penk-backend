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
	UpdateMetricsList  *graphql.Field
	DeleteCustomMetric *graphql.Field
	ResetCustomMetric  *graphql.Field

	CreateMetricProperty *graphql.Field
	UpdateMetricProperty *graphql.Field
	DeleteMetricProperty *graphql.Field
}

func InitCharacterMutation(r *CharactersResolver) *CharactersMutation {
	return &CharactersMutation{
		CreateCharacter: &graphql.Field{
			Type:        CharacterType,
			Description: "Create a character",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(characterInputType),
				},
			},
			Resolve: r.CreateCharacter,
		},
		UpdateCharacter: &graphql.Field{
			Type:        CharacterType,
			Description: "Update a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(characterInputType),
				},
			},
			Resolve: r.UpdateCharacter,
		},
		DeleteCharacter: &graphql.Field{
			Type:        CharacterType,
			Description: "Delete a character",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.DeleteCharacter,
		},
		ResetCharacter: &graphql.Field{
			Type:        CharacterType,
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
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(customMetricInputType),
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
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(customMetricInputType),
				},
			},
			Resolve: r.UpdateCustomMetric,
		},
		UpdateMetricsList: &graphql.Field{
			Type:        graphql.NewList(customMetricType),
			Description: "Update a Custom Metric",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"characterID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewList(customMetricInputType),
				},
			},
			Resolve: r.UpdateMetricList,
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
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(metricPropertyInputType),
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
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(metricPropertyInputType),
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
