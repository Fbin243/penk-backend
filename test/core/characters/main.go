package characters

import (
	"context"
	"log"

	"tenkhours/test/common"
)

type CreateCharacterStage struct {
	common.Metadata
}

func (s CreateCharacterStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Name)

	user, ok := (*ctx).Value(common.User).(map[string]interface{})
	if !ok {
		return common.ErrNotFoundInContext(common.User)
	}

	assertion := CreateCharacterAssertion(user["id"])
	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     CreateCharacterQuery,
			Variables: CreateCharacterVariable,
			Assertion: assertion.End(),
		})
}

type UpdateCharacterStage struct {
	common.Metadata
	CharacterKey common.ContextKey
}

func (s UpdateCharacterStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Name)

	character, ok := (*ctx).Value(s.CharacterKey).(map[string]interface{})
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	assertion := UpdateCharacterAssertion(character["id"])
	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UpdateCharacterQuery,
			Variables: UpdateCharacterVariable(character["id"]),
			Assertion: assertion.End(),
		})
}

type DeleteCharacterStage struct {
	common.Metadata
	CharacterKey common.ContextKey
}

func (s DeleteCharacterStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Name)

	character, ok := (*ctx).Value(s.CharacterKey).(map[string]interface{})
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	assertion := common.AssertionSuccess
	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query: DeleteCharacterQuery,
			Variables: map[string]interface{}{
				"id": character["id"],
			},
			Assertion: assertion.End(),
		})
}
