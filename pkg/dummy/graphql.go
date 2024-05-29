package dummy

import (
	"fmt"
	"slices"

	"github.com/graphql-go/graphql"
)

var localDummies []*DummyData = createDummyArrays()

var dummyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Dummy",
	Fields: graphql.Fields{
		"timestamp": &graphql.Field{
			Type: graphql.DateTime,
		},
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"message": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"dummy": &graphql.Field{
			Type:        dummyType,
			Description: "Get a dummy",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)

				i := slices.IndexFunc(localDummies, func(dd *DummyData) bool {
					return dd.ID == id
				})

				if i == -1 {
					return nil, fmt.Errorf("dummy not found")
				}

				return localDummies[i], nil
			},
		},
		"dummies": &graphql.Field{
			Type:        graphql.NewList(dummyType),
			Description: "Get all dummies",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return localDummies, nil
			},
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"updateDummy": &graphql.Field{
			Type:        dummyType,
			Description: "Update a dummy by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"message": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)

				message := params.Args["message"].(string)

				i := slices.IndexFunc(localDummies, func(dd *DummyData) bool {
					return dd.ID == id
				})

				if i == -1 {
					return nil, fmt.Errorf("dummy not found")
				}

				localDummies[i].Message = message

				return localDummies[i], nil
			},
		},
	},
})

var DummySchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
