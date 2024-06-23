package characters

import (
	"github.com/graphql-go/graphql"
)

var CreateCharacter = graphql.Field{
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
}

var UpdateCharacter = graphql.Field{
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
		"totalFocusTime": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"customMetrics": &graphql.ArgumentConfig{
			Type: graphql.NewList(newCustomMetricInput),
		},
	},
	Resolve: updateCharacter,
}

var DeleteCharacter = graphql.Field{
	Type:        graphql.Boolean,
	Description: "Delete a character",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: deleteCharacter,
}

var ResetCharacter = graphql.Field{
	Type:        graphql.Boolean,
	Description: "Reset a character",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: resetCharacter,
}

var CreateCustomMetric = graphql.Field{
	Type:        customMetricsType,
	Description: "Create a Custom Metrics",
	Args: graphql.FieldConfigArgument{
		"CharacterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"Name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"Description": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"style": &graphql.ArgumentConfig{
			Type: styleType,
		},
	},
	Resolve: createCustomMetric,
}

var UpdateCustomMetric = graphql.Field{
	Type:        customMetricsType,
	Description: "Create a Custom Metrics",
	Args: graphql.FieldConfigArgument{
		"ID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"CharacterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"Name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"Description": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"style": &graphql.ArgumentConfig{
			Type: styleType,
		},
	},
	Resolve: updateCustomMetric,
}

var ResetCustomMetric = graphql.Field{
	Type:        graphql.Boolean,
	Description: "Reset a custom metric",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"CharacterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: resetCustomMetric,
}

var DeleteCustomMetric = graphql.Field{
	Type:        graphql.Boolean,
	Description: "Delete a custom metric",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"CharacterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: deleteCustomMetric,
}

var CreateMetricProperty = graphql.Field{
	Type:        metricProperty,
	Description: "Create a Custom Metrics Property",
	Args: graphql.FieldConfigArgument{
		"CharacterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"MetricID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"Name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"Type": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"Value": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"Unit": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: createCustomMetricProperty,
}

var UpdateMetricProperty = graphql.Field{
	Type:        metricProperty,
	Description: "Update a Custom Metrics Property",
	Args: graphql.FieldConfigArgument{
		"ID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"CharacterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"MetricID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"Name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"Type": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"Value": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"Unit": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: updateCustomMetricProperty,
}

var DeleteCustomMetricProperty = graphql.Field{
	Type:        graphql.Boolean,
	Description: "Delete a custom metric property",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"CharacterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"MetricID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: deleteCustomMetricProperty,
}

var ResetCustomMetricProperty = graphql.Field{
	Type:        graphql.Boolean,
	Description: "Reset a custom metric property",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"CharacterID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"MetricID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: resetCustomMetricProperty,
}
