package core

import (
	"context"
	"log"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/tidwall/gjson"
)

type CreateCustomMetricStage struct {
	common.Metadata
	CharacterKey common.ContextKey // The key to the character in the context
}

func (s CreateCustomMetricStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	variables := map[string]interface{}{
		"characterID": gjson.Get(character, "id").Value(),
		"name":        "Metric name",
		"description": "This is the custom metric description",
		"style": map[string]interface{}{
			"color": "#000000",
			"icon":  "icon.png",
		},
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Present("$.data.createCustomMetric.id").
		Equal("$.data.createCustomMetric.name", variables["name"]).
		Equal("$.data.createCustomMetric.description", variables["description"]).
		Equal("$.data.createCustomMetric.style", variables["style"]).
		Equal("$.data.createCustomMetric.time", float64(0)).
		Equal("$.data.createCustomMetric.limitedPropertyNumber", float64(2))

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     CreateCustomMetricQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type UpdateCustomMetricStage struct {
	common.Metadata
	CharacterKey    common.ContextKey
	CustomMetricKey common.ContextKey
}

func (s UpdateCustomMetricStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CustomMetricKey)
	}

	variables := map[string]interface{}{
		"id":          gjson.Get(customMetric, "id").Value(),
		"characterID": gjson.Get(character, "id").Value(),
		"name":        "Update metric name",
		"description": "This is the updated custom metric description",
		"style": map[string]interface{}{
			"color": "#123456",
			"icon":  "update_icon.png",
		},
		"properties": []interface{}{
			map[string]interface{}{
				"name":  "property 1",
				"type":  "type 1",
				"value": "value 1",
				"unit":  "unit 1",
			},
			map[string]interface{}{
				"name":  "property 1",
				"type":  "type 1",
				"value": "value 1",
				"unit":  "unit 1",
			},
		},
	}

	assertion := jsonpath.Chain().
		Present("data.updateCustomMetric.id").
		Equal("data.updateCustomMetric.name", variables["name"]).
		Equal("data.updateCustomMetric.description", variables["description"]).
		Equal("data.updateCustomMetric.style", variables["style"]).
		Equal("data.updateCustomMetric.properties", variables["properties"]).
		NotEqual("data.updateCustomMetric.time", nil).
		NotEqual("data.updateCustomMetric.limitedPropertyNumber", nil)

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UpdateCustomMetricQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type DeleteCustomMetricStage struct {
	common.Metadata
	CharacterKey    common.ContextKey
	CustomMetricKey common.ContextKey
}

func (s DeleteCustomMetricStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CustomMetricKey)
	}

	variables := map[string]interface{}{
		"id":          gjson.Get(customMetric, "id").Value(),
		"characterID": gjson.Get(character, "id").Value(),
	}

	assertion := common.AssertionSuccess

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UpdateCustomMetricQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type ResetCustomMetricStage struct {
	common.Metadata
	CharacterKey    common.ContextKey
	CustomMetricKey common.ContextKey
}

func (s ResetCustomMetricStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CustomMetricKey)
	}

	variables := map[string]interface{}{
		"id":          gjson.Get(customMetric, "id").Value(),
		"characterID": gjson.Get(character, "id").Value(),
	}

	afterReset := map[string]interface{}{
		"description": nil,
		"style": map[string]interface{}{
			"color": nil,
			"icon":  nil,
		},
		"time":       float64(0),
		"properties": []interface{}{},
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.resetCustomMetric.id", variables["id"]).
		Equal("$.data.resetCustomMetric.description", afterReset["description"]).
		Equal("$.data.resetCustomMetric.style", afterReset["style"]).
		Equal("$.data.resetCustomMetric.time", afterReset["time"])

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     ResetCustomMetricQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}
