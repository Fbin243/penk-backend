package users

import (
	"context"

	"tenkhours/pineline"
	"tenkhours/test/common"
)

func GetUser(expectError bool) pineline.Stage {
	return common.QueryGraphQL(func(ctx *context.Context) (*common.QueryParams, error) {
		assertion := NewUserAssertion
		if expectError {
			assertion = common.AssertionError
		}

		return &common.QueryParams{
			Query:     UserQuery,
			Assertion: assertion.End(),
		}, nil
	})
}

func UpdateUser(expectError bool) pineline.Stage {
	return common.QueryGraphQL(func(ctx *context.Context) (*common.QueryParams, error) {
		assertion := UpdateAccountAssertion
		if expectError {
			assertion = common.AssertionError
		}

		return &common.QueryParams{
			Query:     UpdateAccountQuery,
			Variables: UpdateAccountVariable,
			Assertion: assertion.End(),
		}, nil
	})
}
