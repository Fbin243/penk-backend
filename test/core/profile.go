package core

import (
	"context"
	"log"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

type UpsertProfileStage struct {
	common.Metadata
	common.Case
}

func (s UpsertProfileStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	assertion := jsonpath.Chain().NotPresent("$.errors")
	query := ProfileQuery
	variables := map[string]interface{}{}

	switch s.Case {
	case common.GetProfile:
		assertion = assertion.
			NotEqual("$.data.profile.id", "").
			NotEqual("$.data.profile.firebaseUID", "").
			Equal("$.data.profile.characters", []interface{}{}).
			Equal("$.data.profile.currentCharacterID", nil)

	case common.UpdateProfile:
		profileInput := map[string]interface{}{
			"name":               "Profile name",
			"imageURL":           "https://image.com",
			"currentCharacterID": "675d4ce886be5c6dd755542f",
		}

		variables["input"] = profileInput
		query = UpdateProfileQuery
		assertion = assertion.
			NotEqual("$.data.updateProfile.id", "").
			Equal("$.data.updateProfile.name", profileInput["name"]).
			Equal("$.data.updateProfile.imageURL", profileInput["imageURL"]).
			Equal("$.data.updateProfile.currentCharacterID", profileInput["currentCharacterID"])
	}

	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     query,
			Variables: variables,
			Assertion: []common.Assertion{assertion.End()},
		})
}
