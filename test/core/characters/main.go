package characters

import (
	"context"
	"log"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/tidwall/gjson"
)

type CreateCharacterStage struct {
	common.Metadata
}

func (s CreateCharacterStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	user, ok := (*ctx).Value(common.User).(string)
	if !ok {
		return common.ErrNotFoundInContext(common.User)
	}

	variables := map[string]interface{}{
		"name":   "Character name",
		"gender": false,
		"avatar": "avatar.png",
		"tags":   []interface{}{"#Tag1", "#Tag2"},
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Present("$.data.createCharacter.id").
		Equal("$.data.createCharacter.name", variables["name"]).
		Equal("$.data.createCharacter.avatar", variables["avatar"]).
		Equal("$.data.createCharacter.gender", variables["gender"]).
		Equal("$.data.createCharacter.tags", variables["tags"]).
		Equal("$.data.createCharacter.limitedMetricNumber", float64(2)).
		Equal("$.data.createCharacter.totalFocusedTime", float64(0)).
		Equal("$.data.createCharacter.userID", gjson.Get(user, "id").Value()).
		Equal("$.data.createCharacter.customMetrics", []interface{}{})

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.CoreUrl,
			Query:     CreateCharacterQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type UpdateCharacterStage struct {
	common.Metadata
	CharacterKey common.ContextKey
}

func (s UpdateCharacterStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	variables := map[string]interface{}{
		"id":     gjson.Get(character, "id").Value(),
		"name":   "Update name",
		"gender": true,
		"avatar": "update-avatar.png",
		"tags":   []interface{}{"#update_tag_1", "#update_tag_2"},
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.updateCharacter.id", variables["id"]).
		Equal("$.data.updateCharacter.name", variables["name"]).
		Equal("$.data.updateCharacter.avatar", variables["avatar"]).
		Equal("$.data.updateCharacter.tags", variables["tags"])

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.CoreUrl,
			Query:     UpdateCharacterQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type DeleteCharacterStage struct {
	common.Metadata
	CharacterKey common.ContextKey
}

func (s DeleteCharacterStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	variables := map[string]interface{}{
		"id": gjson.Get(character, "id").Value(),
	}

	assertion := common.AssertionSuccess
	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.CoreUrl,
			Query:     DeleteCharacterQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}
