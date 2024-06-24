package characters

import (
	"github.com/graphql-go/graphql"
)

var styleType = graphql.NewObject(graphql.ObjectConfig{
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

var metricProperty = graphql.NewObject(graphql.ObjectConfig{
	Name: "MetricProperty",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
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

var customMetricsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CustomMetrics",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
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
			Type: styleType,
		},
		"properties": &graphql.Field{
			Type: graphql.NewList(metricProperty),
		},
		"limitedProperties": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var characterType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Character",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"userID": &graphql.Field{
			Type: graphql.String,
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
			Type: graphql.NewList(customMetricsType),
		},
		"limitedCustomMetrics": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var styleTypeInput = graphql.NewInputObject(graphql.InputObjectConfig{
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

var metricPropertyInputType = graphql.NewInputObject(graphql.InputObjectConfig{
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
