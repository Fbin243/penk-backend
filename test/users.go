package test

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func registerNewUser(t *testing.T, ctx *TestContext) {
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(`mutation { registerAccount }`).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}

func getUserInfo(t *testing.T, ctx *TestContext) {
	// Get the user by ID -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(`query { 
			user {createdAt currentCharacterID email firebaseUID id imageURL name updatedAt}
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}
