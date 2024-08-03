package test

import (
	"context"

	"tenkhours/pineline"
	"tenkhours/test/assertions"
	"tenkhours/test/query"
	"tenkhours/test/variables"
)

/**
 * * Create character
 */
func createCharacter(expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		user, ok := (*ctx).Value(User).(map[string]interface{})
		if !ok {
			return nil, ErrNotFoundInContext(User)
		}

		assertion := assertions.CreateCharacter(user["id"])
		if expectError {
			assertion = assertions.AssertionError
		}

		return &QueryParams{
			Query:     query.CreateCharacter,
			Variables: variables.CreateCharacter,
			Assertion: assertion.End(),
		}, nil
	})
}

/**
 * * Read character
 */

/**
 * * Update character
 */
func updateCharacter(characterKey ContextKey, expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		character, ok := (*ctx).Value(characterKey).(map[string]interface{})
		if !ok {
			return nil, ErrNotFoundInContext(characterKey)
		}

		assertion := assertions.UpdateCharacter(character["id"])
		if expectError {
			assertion = assertions.AssertionError
		}

		return &QueryParams{
			Query:     query.UpdateCharacter,
			Variables: variables.UpdateCharacter(character["id"]),
			Assertion: assertion.End(),
		}, nil
	})
}

// /**
//   - * Delete character
//     */
func deleteCharacter(characterKey ContextKey, expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		character, ok := (*ctx).Value(characterKey).(map[string]interface{})
		if !ok {
			return nil, ErrNotFoundInContext(characterKey)
		}

		assertion := assertions.AssertionSuccess
		if expectError {
			assertion = assertions.AssertionError
		}

		return &QueryParams{
			Query: query.DeleteCharacter,
			Variables: map[string]interface{}{
				"id": character["id"],
			},
			Assertion: assertion.End(),
		}, nil
	})
}
