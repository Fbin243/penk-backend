package test

import (
	"context"

	"tenkhours/pineline"
	"tenkhours/test/assertions"
	"tenkhours/test/query"
	"tenkhours/test/variables"
)

func createCustomMetric(characterKey ContextKey, expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		character, ok := (*ctx).Value(characterKey).(map[string]interface{})
		if !ok {
			return nil, ErrNotFoundInContext(characterKey)
		}

		assertion := assertions.CreateCustomMetric
		if expectError {
			assertion = assertions.AssertionError
		}

		return &QueryParams{
			Query:     query.CreateCustomMetric,
			Variables: variables.CreateCustomMetric(character["id"]),
			Assertion: assertion.End(),
		}, nil
	})
}
