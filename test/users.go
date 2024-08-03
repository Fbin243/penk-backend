package test

import (
	"context"

	"tenkhours/pineline"
	"tenkhours/test/assertions"
	"tenkhours/test/query"
	"tenkhours/test/variables"
)

func getUser(expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		assertion := assertions.AssertionSuccess
		if expectError {
			assertion = assertions.AssertionError
		}

		return &QueryParams{
			Query:     query.User,
			Assertion: assertion.End(),
		}, nil
	})
}

func updateUser(expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		variables := variables.UpdateAccount

		assertion := assertions.UpdateAccount
		if expectError {
			assertion = assertions.AssertionError
		}

		return &QueryParams{
			Query:     query.UpdateAccount,
			Variables: variables,
			Assertion: assertion.End(),
		}, nil
	})
}
