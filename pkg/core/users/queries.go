package users

import "github.com/graphql-go/graphql"

var userQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserQuery",
	Fields: graphql.Fields{
		"getUser": &graphql.Field{
			Type:        userType,
			Description: "Get a user by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
		},
	},
})
