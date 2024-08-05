package characters

import (
	"context"

	"tenkhours/pineline"
	"tenkhours/test/common"
)

/**
 * * Create a new character
 */
func CreateCharacter(expectError bool) pineline.Stage {
	return common.QueryGraphQL(func(ctx *context.Context) (*common.QueryParams, error) {
		user, ok := (*ctx).Value(common.User).(map[string]interface{})
		if !ok {
			return nil, common.ErrNotFoundInContext(common.User)
		}

		assertion := CreateCharacterAssertion(user["id"])
		if expectError {
			assertion = common.AssertionError
		}

		return &common.QueryParams{
			Query:     CreateCharacterQuery,
			Variables: CreateCharacterVariable,
			Assertion: assertion.End(),
		}, nil
	})
}

/**
 * * Update character
 */
func UpdateCharacter(characterKey common.ContextKey, expectError bool) pineline.Stage {
	return common.QueryGraphQL(func(ctx *context.Context) (*common.QueryParams, error) {
		character, ok := (*ctx).Value(characterKey).(map[string]interface{})
		if !ok {
			return nil, common.ErrNotFoundInContext(characterKey)
		}

		assertion := UpdateCharacterAssertion(character["id"])
		if expectError {
			assertion = common.AssertionError
		}

		return &common.QueryParams{
			Query:     UpdateCharacterQuery,
			Variables: UpdateCharacterVariable(character["id"]),
			Assertion: assertion.End(),
		}, nil
	})
}

/**
* * Delete character
 */
func DeleteCharacter(characterKey common.ContextKey, expectError bool) pineline.Stage {
	return common.QueryGraphQL(func(ctx *context.Context) (*common.QueryParams, error) {
		character, ok := (*ctx).Value(characterKey).(map[string]interface{})
		if !ok {
			return nil, common.ErrNotFoundInContext(characterKey)
		}

		assertion := common.AssertionSuccess
		if expectError {
			assertion = common.AssertionError
		}

		return &common.QueryParams{
			Query: DeleteCharacterQuery,
			Variables: map[string]interface{}{
				"id": character["id"],
			},
			Assertion: assertion.End(),
		}, nil
	})
}
