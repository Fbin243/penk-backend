package characters

import (
	"fmt"

	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
)

var metricStyleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StyleType",
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

				return nil, fmt.Errorf("failed to convert property ObjectID to Hex")
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

var customMetricType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CustomMetrics",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if metric, ok := p.Source.(coredb.CustomMetric); ok {
					return metric.ID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to convert custom metric ObjectID to Hex")
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

var characterType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Character",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if character, ok := p.Source.(coredb.Character); ok {
					return character.ID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to convert character ObjectID to Hex")
			},
		},
		"userID": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if character, ok := p.Source.(coredb.Character); ok {
					return character.UserID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to convert user ObjectID to Hex")
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"tags": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"totalFocusTime": &graphql.Field{
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

var metricStyleTypeInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "StyleTypeInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"color": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"icon": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

var metricPropertyTypeInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MetricPropertyInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"type": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"value": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"unit": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})
