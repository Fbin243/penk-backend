package users

import (
	"context"
	"log"

	"tenkhours/test/common"
)

type GetUserStage struct {
	common.Metadata
}

func (s GetUserStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Name)

	assertion := NewUserAssertion
	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UserQuery,
			Assertion: assertion.End(),
		})
}

type UpdateUserStage struct {
	common.Metadata
}

func (s UpdateUserStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Name)

	assertion := UpdateAccountAssertion
	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UpdateAccountQuery,
			Variables: UpdateAccountVariable,
			Assertion: assertion.End(),
		})
}
