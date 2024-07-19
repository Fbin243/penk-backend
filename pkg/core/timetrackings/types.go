package timetrackings

import (
	"fmt"

	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

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

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"characterID": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(coredb.TimeTracking); ok {
					return timeTracking.CharacterID.Hex(), nil
				}

				return nil, utils.ErrorConvertOIDToHex
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

				return nil, utils.ErrorConvertOIDToHex
			},
		},
		"startTime": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(coredb.TimeTracking); ok {
					return timeTracking.StartTime.Unix(), nil
				}

				return nil, fmt.Errorf("failed to convert time tracking start time to Unix")
			},
		},
		"endTime": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(coredb.TimeTracking); ok {
					if timeTracking.EndTime.IsZero() {
						return nil, nil
					}

					return timeTracking.EndTime.Unix(), nil
				}

				return nil, fmt.Errorf("failed to convert time tracking end time to Unix")
			},
		},
	},
})
