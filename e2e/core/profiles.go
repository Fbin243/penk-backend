package core

import (
	"context"
	"log"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

type GetProfileStage struct {
	common.Metadata
	CreateNewProfile bool
}

func (s GetProfileStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	assertion := common.AssertionSuccess

	if s.CreateNewProfile {
		assertion = jsonpath.Chain().NotPresent("$.errors").
			Present("$.data.profile.id").
			Present("$.data.profile.firebaseUID").
			Equal("$.data.profile.availableSnapshots", float64(2)).
			Equal("$.data.profile.autoSnapshot", true).
			Equal("$.data.profile.characters", []interface{}{}).
			Equal("$.data.profile.currentCharacterID", nil)
	}

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.CoreUrl,
			Query:     ProfileQuery,
			Assertion: assertion.End(),
		})
}

type UpdateProfileStage struct {
	common.Metadata
}

func (s UpdateProfileStage) Exec(ctx *context.Context) error {
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
