package core

import (
	"tenkhours/pkg/core/characters"

	"github.com/graphql-go/graphql"
)

// var rootQuery = graphql.NewObject(graphql.ObjectConfig{
// 	Name: "RootQuery",
// 	Fields: graphql.Fields{
// 		"user": &users.User,
// 	},
// })

// var rootMutation = graphql.NewObject(graphql.ObjectConfig{
// 	Name: "RootMutation",
// 	Fields: graphql.Fields{
// 		"registerAccount": &users.RegisterAccount,
// 	},
// })

// var CoreSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
// 	Query:    rootQuery,
// 	Mutation: rootMutation,
// })

var rootQuery = characters.RootQuery
var rootMutation = characters.RootMutation

var CoreSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
