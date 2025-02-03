package core

import (
	"context"
	"log"

	"tenkhours/pkg/utils"
	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/tidwall/gjson"
)

type UpsertCharacterStage struct {
	common.Metadata
	common.Case
	CharacterKey    common.ContextKey
	NumberOfMetrics int
}

func (s UpsertCharacterStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	profile, ok := (*ctx).Value(common.Profile).(string)
	if !ok {
		return common.ErrNotFoundInContext("Profile")
	}

	characterInput := map[string]interface{}{
		"name":   "Character name",
		"gender": false,
		"tags":   []interface{}{"#Tag1", "#Tag2"},
	}

	metricInput := map[string]interface{}{
		"name":        "Metric name",
		"description": "This is the metric description",
		"style": map[string]interface{}{
			"color": "#000000",
			"icon":  "icon.png",
		},
	}

	assertion := jsonpath.Chain().NotPresent("$.errors")
	assertions := []common.Assertion{}
	switch s.Case {
	case common.CreateCharacter:
		assertion.
			NotEqual("$.data.upsertCharacter.id", "").
			Equal("$.data.upsertCharacter.name", characterInput["name"]).
			Equal("$.data.upsertCharacter.gender", characterInput["gender"]).
			Equal("$.data.upsertCharacter.tags", characterInput["tags"]).
			Equal("$.data.upsertCharacter.limitedMetricNumber", float64(utils.LimitedMetricNumber)).
			Equal("$.data.upsertCharacter.totalFocusedTime", float64(0)).
			Equal("$.data.upsertCharacter.profileID", gjson.Get(profile, "id").Value()).
			Equal("$.data.upsertCharacter.customMetrics", []interface{}{})

	case common.UpdateCharacter:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		characterInput["id"] = gjson.Get(character, "id").Value()
		characterInput["name"] = "Update name"
		characterInput["tags"] = []interface{}{"#update_tag_1", "#update_tag_2"}

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", characterInput["id"]).
			Equal("$.data.upsertCharacter.name", characterInput["name"]).
			Equal("$.data.upsertCharacter.tags", characterInput["tags"])

	case common.CreateMetrics:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		metricInputs := []interface{}{}
		for i := 0; i < s.NumberOfMetrics; i++ {
			metricInputs = append(metricInputs, metricInput)
		}
		characterInput["id"] = gjson.Get(character, "id").Value()
		characterInput["customMetrics"] = metricInputs

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", characterInput["id"]).
			Equal("$.data.upsertCharacter.customMetrics[0].name", metricInput["name"]).
			Equal("$.data.upsertCharacter.customMetrics[0].description", metricInput["description"]).
			Equal("$.data.upsertCharacter.customMetrics[0].style", metricInput["style"].(map[string]interface{})).
			Equal("$.data.upsertCharacter.customMetrics[0].time", float64(0)).
			Equal("$.data.upsertCharacter.customMetrics[0].limitedPropertyNumber", float64(utils.LimitedPropertyNumber))
		assertions = append(assertions, jsonpath.Len("$.data.upsertCharacter.customMetrics", s.NumberOfMetrics))

	case common.UpdateMetrics:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		// Update the first metric
		characterInput["id"] = gjson.Get(character, "id").Value()
		metricInput["id"] = gjson.Get(character, "customMetrics.0.id").Value()
		metricInput["name"] = "Update metric name"
		metricInput["description"] = "Update metric description"
		metricInput["style"] = map[string]interface{}{
			"color": "#FFFFFF",
			"icon":  "update_icon.png",
		}

		metricInputs := []interface{}{metricInput}
		for i := 0; i < s.NumberOfMetrics-1; i++ {
			metricInputs = append(metricInputs, metricInput)
		}

		characterInput["customMetrics"] = metricInputs

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", gjson.Get(character, "id").Value()).
			Equal("$.data.upsertCharacter.customMetrics[0].name", metricInput["name"]).
			Equal("$.data.upsertCharacter.customMetrics[0].description", metricInput["description"]).
			Equal("$.data.upsertCharacter.customMetrics[0].style", metricInput["style"].(map[string]interface{}))
		assertions = append(assertions, jsonpath.Len("$.data.upsertCharacter.customMetrics", s.NumberOfMetrics))

	case common.DeleteMetrics:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		// Remove the first metric
		metricInputs := []interface{}{}
		for i := 0; i < s.NumberOfMetrics-1; i++ {
			metricInputs = append(metricInputs, metricInput)
		}

		characterInput["id"] = gjson.Get(character, "id").Value()
		characterInput["customMetrics"] = metricInputs

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", characterInput["id"])

		assertions = append(assertions, jsonpath.Len("$.data.upsertCharacter.customMetrics", s.NumberOfMetrics-1))
	}

	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UpsertCharacterQuery,
			Variables: map[string]interface{}{"input": characterInput},
			Assertion: append(assertions, assertion.End()),
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
		return common.ErrNotFoundInContext("CharacterKey")
	}

	variables := map[string]interface{}{
		"id": gjson.Get(character, "id").Value(),
	}

	assertion := jsonpath.Chain().NotPresent("$.errors")
	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     DeleteCharacterQuery,
			Variables: variables,
			Assertion: []common.Assertion{assertion.End()},
		})
}
