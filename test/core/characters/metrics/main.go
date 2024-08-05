package metrics

import (
	"context"

	"tenkhours/pineline"
	"tenkhours/test/common"
)

/**
 * * Create a new custom metric
 */
func CreateCustomMetric(characterKey common.ContextKey, expectError bool) pineline.Stage {
	return common.QueryGraphQL(func(ctx *context.Context) (*common.QueryParams, error) {
		character, ok := (*ctx).Value(characterKey).(map[string]interface{})
		if !ok {
			return nil, common.ErrNotFoundInContext(characterKey)
		}

		assertion := CreateCustomMetricAssertion
		if expectError {
			assertion = common.AssertionError
		}

		return &common.QueryParams{
			Query:     CreateCustomMetricQuery,
			Variables: CreateCustomMetricVariable(character["id"]),
			Assertion: assertion.End(),
		}, nil
	})
}
