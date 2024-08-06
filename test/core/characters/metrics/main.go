package metrics

import (
	"context"
	"log"

	"tenkhours/test/common"
)

type CreateCustomMetricStage struct {
	common.Metadata
	CharacterKey common.ContextKey // The key to the character in the context
}

func (s CreateCustomMetricStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Name)
	character, ok := (*ctx).Value(s.CharacterKey).(map[string]interface{})
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	assertion := CreateCustomMetricAssertion
	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     CreateCustomMetricQuery,
			Variables: CreateCustomMetricVariable(character["id"]),
			Assertion: assertion.End(),
		})
}
