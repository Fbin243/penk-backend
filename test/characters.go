package test

import (
	"context"

	"tenkhours/pineline"
)

/**
 * * Create character
 */
func createNewCharacter(expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		userID, ok := (*ctx).Value(UserID).(string)
		if !ok {
			return nil, ErrNotFoundInContext(UserID)
		}

		variables := map[string]interface{}{
			"name":   "Character name",
			"gender": false,
			"avatar": "avatar.png",
			"tags":   []interface{}{"#Tag1", "#Tag2"},
		}

		assertionChain := assertionChainSuccess.
			Present("$.data.createCharacter.id").
			Equal("$.data.createCharacter.name", variables["name"]).
			Equal("$.data.createCharacter.avatar", variables["avatar"]).
			Equal("$.data.createCharacter.gender", variables["gender"]).
			Equal("$.data.createCharacter.tags", variables["tags"]).
			Equal("$.data.createCharacter.limitedMetricNumber", float64(2)).
			Equal("$.data.createCharacter.totalFocusedTime", float64(0)).
			Equal("$.data.createCharacter.userID", userID).End()
		if expectError {
			assertionChain = assertionChainError.End()
		}

		return &QueryParams{
			Query: `
			mutation CreateCharacter($name: String!, $gender: Boolean, $avatar: String!, $tags: [String]) {
				createCharacter(
					input: { 
						name: $name, 
						gender: $gender, 
						avatar: $avatar, 
						tags: $tags 
					}
				) {
					avatar
					gender
					id
					limitedMetricNumber
					name
					tags
					totalFocusedTime
					userID
				}
			}`,
			Variables:      variables,
			AssertionChain: assertionChain,
		}, nil
	})
}

/**
 * * Read character
 */

/**
 * * Update character
 */
func updateCharacter(expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		characterID, ok := (*ctx).Value(CharacterID).(string)
		if !ok {
			return nil, ErrNotFoundInContext(CharacterID)
		}

		variables := map[string]interface{}{
			"id":     characterID,
			"name":   "Update name",
			"gender": true,
			"avatar": "update-avatar.png",
			"tags":   []interface{}{"#update_tag_1", "#update_tag_2"},
		}

		assertionChain := assertionChainSuccess.
			Equal("$.data.updateCharacter.name", variables["name"]).
			Equal("$.data.updateCharacter.avatar", variables["avatar"]).
			Equal("$.data.updateCharacter.tags", variables["tags"]).
			End()

		if expectError {
			assertionChain = assertionChainError.End()
		}

		return &QueryParams{
			Query: `
			mutation UpdateCharacter($id: String!, $avatar: String, $gender: Boolean, $name: String, $tags: [String]) {
				updateCharacter(
					id: $id
					input: { 
						avatar: $avatar, 
						gender: $gender, 
						name: $name
						tags: $tags, 
					}
				) {
					avatar
					gender
					name
					tags
				}
			}`,
			Variables:      variables,
			AssertionChain: assertionChain,
		}, nil
	})
}

/**
 * * Delete character
 */
func deleteCharacter(expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		characterID, ok := (*ctx).Value(CharacterID).(string)
		if !ok {
			return nil, ErrNotFoundInContext(CharacterID)
		}

		assertionChain := assertionChainSuccess.End()
		if expectError {
			assertionChain = assertionChainError.End()
		}

		return &QueryParams{
			Query: `
				mutation DeleteCharacter($id: String!) {
					deleteCharacter(id: $id) {
					avatar
					gender
					id
					limitedMetricNumber
					name
					tags
					totalFocusedTime
					userID
				}
			}`,
			Variables: map[string]interface{}{
				"id": characterID,
			},
			AssertionChain: assertionChain,
		}, nil
	})
}
