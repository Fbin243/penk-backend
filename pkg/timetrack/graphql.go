package timetrack

import (
    "github.com/graphql-go/graphql"
)

var TimerType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Timer",
        Fields: graphql.Fields{
            "start": &graphql.Field{
                Type: graphql.DateTime,
            },
            "end": &graphql.Field{
                Type: graphql.DateTime,
            },
        },
    },
)

var StartTimerMutation = &graphql.Field{
    Type:        TimerType,
    Description: "Start a new timer!",
    Resolve: func(params graphql.ResolveParams) (interface{}, error) {
        timer := &Timer{}
        timer.Start()
        return timer, nil
    },
}