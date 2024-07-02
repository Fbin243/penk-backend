package test

import (
	"fmt"
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
		Assert(func(r1 *http.Response, r2 *http.Request) error {
			data, err := decodeResponseData(r1)
			if err != nil {
				return err
			}

			if idUser, ok := data["registerAccount"].(string); ok {
				ctx.IdUser = idUser
				return nil
			}

			return fmt.Errorf("failed to register a new user")
		}).
		End()
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
		Assert(logResponseData).
		End()
}
