package test

import (
	"log"
	"net/http"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var response = &Map{}

func getUserInfo(ctx *Context) error {
	testingT, ok := (*ctx)["testingT"].(apitest.TestingT)
	if !ok {
		return ErrNotFoundInContext
	}

	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(`query { 
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
		}`).
		Expect(testingT).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&response)

	response.log()
	(*ctx)["userID"] = response.getFieldValue("data.user.id")
	log.Print((*ctx)["userID"])

	return nil
}

func updateUser(ctx *Context) error {
	testingT, ok := (*ctx)["testingT"].(apitest.TestingT)
	if !ok {
		return ErrNotFoundInContext
	}

	query := `mutation UpdateAccount($name: String, $imageURL: String, $currentCharacterID: String, $autoSnapshot: Boolean) {
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
	}`

	variables := Map{
		"name":               "Update name",
		"imageURL":           "update.png",
		"currentCharacterID": "669a2bbc53e6629a2931e1be",
		"autoSnapshot":       false,
	}

	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(query, variables).
		Expect(testingT).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		Assert(jsonpath.Equal("$.data.updateAccount.name", variables["name"])).
		Assert(jsonpath.Equal("$.data.updateAccount.imageURL", variables["imageURL"])).
		Assert(jsonpath.Equal("$.data.updateAccount.currentCharacterID", variables["currentCharacterID"])).
		Assert(jsonpath.Equal("$.data.updateAccount.autoSnapshot", variables["autoSnapshot"])).
		End().JSON(&response)

	response.log()
	return nil
}
