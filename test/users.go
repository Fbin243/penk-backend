package test

import (
	"fmt"
	"net/http"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var response = &Map{}

func getUserInfo(ctx *Context) error {
	testingT, ok := (*ctx)["testingT"].(apitest.TestingT)
	if !ok {
		return ErrNotFoundTestingT
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
	return nil
}

func updateUser(ctx *Context) error {
	testingT, ok := (*ctx)["testingT"].(apitest.TestingT)
	if !ok {
		return ErrNotFoundTestingT
	}

	updateInfo := Map{
		"name":               "Update name",
		"imageURL":           "update.png",
		"currentCharacterID": "669a2bbc53e6629a2931e1be",
		"autoSnapshot":       false,
	}

	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation {
			updateAccount(
				input: {
					name: "%s"
					imageURL: "%s"
					currentCharacterID: "%s"
					autoSnapshot: %t
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
		}`, updateInfo["name"], updateInfo["imageURL"], updateInfo["currentCharacterID"], updateInfo["autoSnapshot"])).
		Expect(testingT).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		Assert(jsonpath.Equal("$.data.updateAccount.name", updateInfo["name"])).
		Assert(jsonpath.Equal("$.data.updateAccount.imageURL", updateInfo["imageURL"])).
		Assert(jsonpath.Equal("$.data.updateAccount.currentCharacterID", updateInfo["currentCharacterID"])).
		Assert(jsonpath.Equal("$.data.updateAccount.autoSnapshot", updateInfo["autoSnapshot"])).
		End().JSON(&response)

	response.log()
	return nil
}
