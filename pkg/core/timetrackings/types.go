package timetrackings

import (
	"github.com/graphql-go/graphql"
)

var timeTrackingType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TimeTracking",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"characterID": &graphql.Field{
			Type: graphql.String,
		},
		"startTime": &graphql.Field{
			Type: graphql.Int,
		},
		"endTime": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
