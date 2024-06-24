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
	Resolve: createCustomMetric,
}

var UpdateCustomMetric = graphql.Field{
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
	Resolve: updateCustomMetric,
}

var ResetCustomMetric = graphql.Field{
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
	Resolve: resetCustomMetric,
}

var DeleteCustomMetric = graphql.Field{
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
	Resolve: deleteCustomMetric,
}
