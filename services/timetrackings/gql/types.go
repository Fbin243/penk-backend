package gql

import (
	"fmt"

	"tenkhours/pkg/db/timetrackingsdb"

	"github.com/graphql-go/graphql"
)

var timeTrackingType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TimeTracking",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(timetrackingsdb.TimeTracking); ok {
					return timeTracking.ID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to resolve TimeTracking id")
			},
		},
		"characterID": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(timetrackingsdb.TimeTracking); ok {
					return timeTracking.CharacterID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to resolve TimeTracking characterID")
			},
		},
		"metricID": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(timetrackingsdb.TimeTracking); ok {
					if timeTracking.CustomMetricID.IsZero() {
						return nil, nil
					}

					return timeTracking.CustomMetricID.Hex(), nil
				}

				return nil, fmt.Errorf("failed to resolve TimeTracking metricID")
			},
		},
		"startTime": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"endTime": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if timeTracking, ok := p.Source.(timetrackingsdb.TimeTracking); ok {
					if timeTracking.EndTime.IsZero() {
						return nil, nil
					}

					return timeTracking.EndTime, nil
				}

				return nil, fmt.Errorf("failed to resolve TimeTracking endTime")
			},
		},
	},
})
