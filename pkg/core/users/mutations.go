package users

import (
	"github.com/graphql-go/graphql"
)

var userMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserMutation",
	Fields: graphql.Fields{
		"registerAccount": &graphql.Field{
			Type:        userType,
			Description: "Register a new account",
			Resolve:     RegisterAccount,
		},
	},
})
