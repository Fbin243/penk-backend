package timetrackings

import (
	"fmt"

	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
)

var timeTrackingType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TimeTracking",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(coredb.TimeTracking); ok {
					return timeTracking.ID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to convert time tracking ObjectID to Hex")
			},
		},
		"characterID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(coredb.TimeTracking); ok {
					return timeTracking.CharacterID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to convert character ObjectID to Hex")
			},
		},
		"customMetricID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(coredb.TimeTracking); ok {
					if timeTracking.CustomMetricID.IsZero() {
						return nil, nil
					}

					return timeTracking.CustomMetricID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to convert custom metric ObjectID to Hex")
			},
		},
		"startTime": &graphql.Field{
			Type: graphql.Int,
		},
		"endTime": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
