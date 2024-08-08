package characters

import (
	"tenkhours/pkg/core/characters"
	"tenkhours/pkg/db/analyticsdb"
	"tenkhours/pkg/utils"

	"github.com/graphql-go/graphql"
)

var snapshotType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Snapshot",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if snapshot, ok := p.Source.(analyticsdb.Snapshot); ok {
					return snapshot.ID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"timestamp": &graphql.Field{
			Type: graphql.DateTime,
		},
		"metadata": &graphql.Field{
			Type: metadataType,
		},
		"character": &graphql.Field{
			Type: characters.CharacterType,
		},
	},
})

var metadataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Metadata",
	Fields: graphql.Fields{
		"userID": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if metadata, ok := p.Source.(analyticsdb.Metadata); ok {
					return metadata.UserID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"characterID": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if snapshot, ok := p.Source.(analyticsdb.Metadata); ok {
					return snapshot.CharacterID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
			},
		},
	},
})
