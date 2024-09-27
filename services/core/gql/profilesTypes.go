package gql

import (
	"fmt"

	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
)

var profileType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Profile",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if profile, ok := p.Source.(coredb.Profile); ok {
					return profile.ID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to resolve Profile id")
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
				if profile, ok := p.Source.(coredb.Profile); ok {
					if profile.CurrentCharacterID.IsZero() {
						return nil, nil
					}

					return profile.CurrentCharacterID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to resolve Profile currentCharacterID")
			},
		},
		"characters": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(CharacterType))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if profile, ok := p.Source.(coredb.Profile); ok {
					// Find characters by profile ID
					charactersRepo := coredb.NewCharactersRepo(db.GetDBManager().DB)
					characters, err := charactersRepo.GetCharactersByProfileID(profile.ID)
					if err != nil {
						return nil, err
					}

					return characters, nil
				}

				return nil, fmt.Errorf("failed to get characters by profile id")
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
var profileInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ProfileInput",
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
