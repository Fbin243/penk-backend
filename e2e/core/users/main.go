package users

import (
	"context"
	"log"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

type GetUserStage struct {
	common.Metadata
	CreateNewUser bool
}

func (s GetUserStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	assertion := common.AssertionSuccess

	if s.CreateNewUser {
		assertion = jsonpath.Chain().NotPresent("$.errors").
			Present("$.data.user.id").
			Present("$.data.user.firebaseUID").
			Equal("$.data.user.availableSnapshots", float64(2)).
			Equal("$.data.user.autoSnapshot", true).
			Equal("$.data.user.characters", []interface{}{}).
			Equal("$.data.user.currentCharacterID", nil)
	}

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.CoreUrl,
			Query:     UserQuery,
			Assertion: assertion.End(),
		})
}

type UpdateUserStage struct {
	common.Metadata
}

func (s UpdateUserStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	variables := map[string]interface{}{
		"name":               "Update name",
		"imageURL":           "update.png",
		"currentCharacterID": "669a2bbc53e6629a2931e1be",
		"autoSnapshot":       false,
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.updateAccount.name", variables["name"]).
		Equal("$.data.updateAccount.imageURL", variables["imageURL"]).
		Equal("$.data.updateAccount.currentCharacterID", variables["currentCharacterID"]).
		Equal("$.data.updateAccount.autoSnapshot", variables["autoSnapshot"])

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.CoreUrl,
			Query:     UpdateAccountQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}
