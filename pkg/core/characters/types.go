package characters

import (
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/graphql-go/graphql"
)

var CharacterType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Character",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if character, ok := p.Source.(coredb.Character); ok {
					if character.ID.IsZero() {
						return nil, nil
					}

					return character.ID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"userID": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if character, ok := p.Source.(coredb.Character); ok {
					if character.ID.IsZero() {
						return nil, nil
					}

					return character.UserID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"gender": &graphql.Field{
			Type: graphql.Boolean,
		},
		"avatar": &graphql.Field{
			Type: graphql.String,
		},
		"tags": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"totalFocusedTime": &graphql.Field{
			Type: graphql.Int,
		},
		"customMetrics": &graphql.Field{
			Type: graphql.NewList(customMetricType),
		},
		"limitedMetricNumber": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var customMetricType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CustomMetric",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if metric, ok := p.Source.(coredb.CustomMetric); ok {
					return metric.ID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"time": &graphql.Field{
			Type: graphql.Int,
		},
		"style": &graphql.Field{
			Type: metricStyleType,
		},
		"properties": &graphql.Field{
			Type: graphql.NewList(metricPropertyType),
		},
		"limitedPropertyNumber": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var metricStyleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MetricStyle",
	Fields: graphql.Fields{
		"color": &graphql.Field{
			Type: graphql.String,
		},
		"icon": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var metricPropertyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MetricProperty",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if property, ok := p.Source.(coredb.MetricProperty); ok {
					return property.ID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"value": &graphql.Field{
			Type: graphql.String,
		},
		"unit": &graphql.Field{
			Type: graphql.String,
		},
	},
})

/**
 * Input types
 */
var characterInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CharacterInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"gender": &graphql.InputObjectFieldConfig{
			Type:        graphql.Boolean,
			Description: "Male is true, Female is false. If not specified, it is false by default",
		},
		"avatar": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "URL or file path of the character's avatar",
		},
		"tags": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewList(graphql.String),
			Description: "List of string tags that describe the character",
		},
	},
})

var customMetricInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CustomMetricInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"style": &graphql.InputObjectFieldConfig{
			Type:        metricStyleInputType,
			Description: "Visual style of the metric that be displayed on screen",
		},
	},
})

var metricStyleInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MetricStyleInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"color": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "Color in Hex format",
		},
		"icon": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "URL or file path of the icon",
		},
	},
})

var metricPropertyInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MetricPropertyInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"type": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "Data type of the property",
		},
		"value": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "Specific value of the property type",
		},
		"unit": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "Unit of the property value",
		},
	},
})
