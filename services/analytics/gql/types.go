package gql

import (
	"fmt"

	"tenkhours/pkg/db/analyticsdb"
	"tenkhours/services/core/gql"

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

				return nil, fmt.Errorf("failed to resolve Snapshot id")
			},
		},
		"timestamp": &graphql.Field{
			Type: graphql.DateTime,
		},
		"character": &graphql.Field{
			Type: gql.CharacterType,
		},
	},
})
