package test

import (
	"context"

	"tenkhours/pineline"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func getUser() pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		return &QueryParams{
			Query: `query { 
				user {
					autoSnapshot
					availableSnapshots
					createdAt
					currentCharacterID
					email
					firebaseUID
					id
					imageURL
					name
					updatedAt
				}
			}`,
			AssertionChain: jsonpath.Chain().NotPresent("$.errors").End(),
		}, nil
	})
}

func updateUser() pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		variables := map[string]interface{}{
			"name":               "Update name",
			"imageURL":           "update.png",
			"currentCharacterID": "669a2bbc53e6629a2931e1be",
			"autoSnapshot":       false,
		}

		return &QueryParams{
			Query: `mutation UpdateAccount($name: String, $imageURL: String, $currentCharacterID: String, $autoSnapshot: Boolean) {
				updateAccount(
					input: {
						name: $name
						imageURL: $imageURL
						currentCharacterID: $currentCharacterID
						autoSnapshot: $autoSnapshot
					}
				) {
					autoSnapshot
					availableSnapshots
					createdAt
					currentCharacterID
					email
					firebaseUID
					id
					imageURL
					name
					updatedAt
				}
			}`,
			Variables: variables,
			AssertionChain: jsonpath.Chain().NotPresent("$.errors").
				Equal("$.data.updateAccount.name", variables["name"]).
				Equal("$.data.updateAccount.imageURL", variables["imageURL"]).
				Equal("$.data.updateAccount.currentCharacterID", variables["currentCharacterID"]).
				Equal("$.data.updateAccount.autoSnapshot", variables["autoSnapshot"]).
				End(),
		}, nil
	})
}
