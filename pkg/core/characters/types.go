package characters

import "github.com/graphql-go/graphql"

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
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "CustomMetrics",
				Fields: graphql.Fields{
					"id": &graphql.Field{
						Type: graphql.String,
					},
					"characterID": &graphql.Field{
						Type: graphql.String,
					},
					"type": &graphql.Field{
						Type: graphql.String,
					},
					"name": &graphql.Field{
						Type: graphql.String,
					},
					"value": &graphql.Field{
						Type: graphql.String,
					},
				},
			})),
		},
	},
})
