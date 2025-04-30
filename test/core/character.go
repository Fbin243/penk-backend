package core

import (
	"context"
	"log"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/tidwall/gjson"
)

type UpsertCharacterStage struct {
	common.Metadata
	common.Case
	CharacterKey common.ContextKey
}

func (s UpsertCharacterStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	profile, ok := (*ctx).Value(common.Profile).(string)
	if !ok {
		return common.ErrNotFoundInContext("Profile")
	}

	metricInput := map[string]interface{}{
		"name":  "Metric name",
		"value": 2.0,
		"unit":  "Metric unit",
	}

	categoryInput := map[string]interface{}{
		"name":        "Category name",
		"description": "Category desc",
		"style": map[string]interface{}{
			"color": "#000000",
			"icon":  "icon.png",
		},
	}

	characterInput := map[string]interface{}{
		"name":   "Character name",
		"gender": false,
		"tags":   []interface{}{"#Tag1", "#Tag2"},
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
			Equal("$.data.upsertCharacter.profileID", gjson.Get(profile, "id").Value()).
			Equal("$.data.upsertCharacter.categories", []interface{}{}).
			Equal("$.data.upsertCharacter.metrics", []interface{}{})

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

	case common.CreateCategories:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		characterInput["id"] = gjson.Get(character, "id").Value()
		characterInput["categories"] = []interface{}{categoryInput, categoryInput, categoryInput}

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", characterInput["id"]).
			Equal("$.data.upsertCharacter.categories[0].name", categoryInput["name"]).
			Equal("$.data.upsertCharacter.categories[0].description", categoryInput["description"]).
			Equal("$.data.upsertCharacter.categories[0].style", categoryInput["style"].(map[string]interface{}))

		assertions = append(assertions, jsonpath.Len("$.data.upsertCharacter.categories", 3))

	case common.UpdateCategories:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		// Update the first category
		characterInput["id"] = gjson.Get(character, "id").Value()
		updateCategoryInput := map[string]interface{}{
			"id":          gjson.Get(character, "categories.0.id").Value(),
			"name":        "Update category name",
			"description": "Update category description",
			"style": map[string]interface{}{
				"color": "#FFFFFF",
				"icon":  "update_icon.png",
			},
		}

		characterInput["categories"] = []interface{}{updateCategoryInput, categoryInput, categoryInput}

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", characterInput["id"]).
			Equal("$.data.upsertCharacter.categories[0].id", updateCategoryInput["id"]).
			Equal("$.data.upsertCharacter.categories[0].name", updateCategoryInput["name"]).
			Equal("$.data.upsertCharacter.categories[0].description", updateCategoryInput["description"]).
			Equal("$.data.upsertCharacter.categories[0].style", updateCategoryInput["style"].(map[string]interface{}))
		assertions = append(assertions, jsonpath.Len("$.data.upsertCharacter.categories", 3))

	case common.DeleteCategories:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		characterInput["id"] = gjson.Get(character, "id").Value()
		characterInput["categories"] = []interface{}{categoryInput, categoryInput}

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", characterInput["id"])

		assertions = append(assertions, jsonpath.Len("$.data.upsertCharacter.categories", 2))

	case common.CreateMetrics:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		characterInput["id"] = gjson.Get(character, "id").Value()
		categoryInputWithMetrics := map[string]interface{}{
			"id":          gjson.Get(character, "categories.0.id").Value(),
			"name":        "Example name",
			"description": "Example desc",
			"style": map[string]interface{}{
				"color": "#000000",
				"icon":  "icon.png",
			},
			"metrics": []interface{}{metricInput, metricInput, metricInput},
		}
		characterInput["categories"] = []interface{}{categoryInputWithMetrics, categoryInput, categoryInput}
		characterInput["metrics"] = []interface{}{metricInput, metricInput, metricInput}

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", characterInput["id"]).
			Equal("$.data.upsertCharacter.categories[0].id", categoryInputWithMetrics["id"]).
			Equal("$.data.upsertCharacter.categories[0].metrics[0].name", metricInput["name"]).
			Equal("$.data.upsertCharacter.categories[0].metrics[0].value", metricInput["value"]).
			Equal("$.data.upsertCharacter.categories[0].metrics[0].unit", metricInput["unit"]).
			Equal("$.data.upsertCharacter.metrics[0].name", metricInput["name"]).
			Equal("$.data.upsertCharacter.metrics[0].value", metricInput["value"]).
			Equal("$.data.upsertCharacter.metrics[0].unit", metricInput["unit"])

		assertions = append(assertions,
			jsonpath.Len("$.data.upsertCharacter.categories[0].metrics", 3),
			jsonpath.Len("$.data.upsertCharacter.metrics", 3))

	case common.UpdateMetrics:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		characterInput["id"] = gjson.Get(character, "id").Value()
		updateMetricInput := map[string]interface{}{
			"id":    gjson.Get(character, "metrics.0.id").Value(),
			"name":  "Update metric name",
			"value": 789.0,
			"unit":  "Update metric unit",
		}

		categoryInputWithUpdateMetrics := map[string]interface{}{
			"id":          gjson.Get(character, "categories.0.id").Value(),
			"name":        "Category name",
			"description": "Category desc",
			"style": map[string]interface{}{
				"color": "#000000",
				"icon":  "icon.png",
			},
			"metrics": []interface{}{updateMetricInput, metricInput, metricInput},
		}

		characterInput["categories"] = []interface{}{categoryInputWithUpdateMetrics, categoryInput, categoryInput}
		characterInput["metrics"] = []interface{}{updateMetricInput, metricInput, metricInput}

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", characterInput["id"]).
			Equal("$.data.upsertCharacter.categories[0].id", categoryInputWithUpdateMetrics["id"]).
			Equal("$.data.upsertCharacter.categories[0].metrics[0].name", updateMetricInput["name"]).
			Equal("$.data.upsertCharacter.categories[0].metrics[0].value", updateMetricInput["value"]).
			Equal("$.data.upsertCharacter.categories[0].metrics[0].unit", updateMetricInput["unit"]).
			Equal("$.data.upsertCharacter.metrics[0].id", updateMetricInput["id"]).
			Equal("$.data.upsertCharacter.metrics[0].name", updateMetricInput["name"]).
			Equal("$.data.upsertCharacter.metrics[0].value", updateMetricInput["value"]).
			Equal("$.data.upsertCharacter.metrics[0].unit", updateMetricInput["unit"])
		assertions = append(assertions,
			jsonpath.Len("$.data.upsertCharacter.categories[0].metrics", 3),
			jsonpath.Len("$.data.upsertCharacter.metrics", 3))

	case common.DeleteMetrics:
		character, ok := (*ctx).Value(s.CharacterKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CharacterKey")
		}

		characterInput["id"] = gjson.Get(character, "id").Value()
		characterInput["metrics"] = []interface{}{metricInput, metricInput}

		assertion = assertion.
			Equal("$.data.upsertCharacter.id", characterInput["id"])

		assertions = append(assertions, jsonpath.Len("$.data.upsertCharacter.metrics", 2))
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
