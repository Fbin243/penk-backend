package users

import (
	"fmt"

	"tenkhours/pkg/core/characters"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(coredb.User); ok {
					return user.ID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"email": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"firebaseUID": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"imageURL": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"currentCharacterID": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(coredb.User); ok {
					if user.CurrentCharacterID.IsZero() {
						return nil, nil
					}

					return user.CurrentCharacterID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"characters": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(characters.CharacterType)),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Source.(coredb.User); ok {
					// Find characters by user ID
					charactersRepo := coredb.NewCharactersRepo(db.GetDBManager().DB)
					characters, err := charactersRepo.GetCharactersByUserID(user.ID)
					if err != nil {
						return nil, err
					}

					return characters, nil
				}

				return nil, fmt.Errorf("failed to get characters by user ID")
			},
		},
		"availableSnapshots": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"autoSnapshot": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
		},
		"createdAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"updatedAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
	},
})

// Input type
var userInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UserInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"imageURL": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "URL of the user's image",
		},
		"currentCharacterID": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "ID of the character is being chosen",
		},
		"autoSnapshot": &graphql.InputObjectFieldConfig{
			Type:        graphql.Boolean,
			Description: "Whether the user has enabled auto snapshot, default is true",
		},
	},
})
