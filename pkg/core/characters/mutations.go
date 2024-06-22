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
