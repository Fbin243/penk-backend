package core

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"tenkhours/test/common"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

type CreateCustomMetricStage struct {
	common.Metadata
	CharacterKey common.ContextKey // The key to the character in the context
	Case         common.CreateCustomMetricCase
}

func (s CreateCustomMetricStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CharacterKey")
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
		NotEqual("$.data.createCustomMetric.id", nil).
		Equal("$.data.createCustomMetric.name", variables["name"]).
		Equal("$.data.createCustomMetric.description", variables["description"]).
		Equal("$.data.createCustomMetric.style", variables["style"]).
		Equal("$.data.createCustomMetric.time", float64(0)).
		Equal("$.data.createCustomMetric.limitedPropertyNumber", float64(2))

	if s.Case == common.CreateMetricWithProperties {
		log.Println("---> Creating custom metric with properties")
		variables["properties"] = []interface{}{
			map[string]interface{}{
				"name":  "property 1",
				"type":  "NUMBER",
				"value": "10",
				"unit":  "Steps",
			},
			map[string]interface{}{
				"name":  "property 2",
				"type":  "STRING",
				"value": "ABCD",
				"unit":  "Characters",
			},
		}

		assertion = assertion.Equal("$.data.createCustomMetric.properties", variables["properties"])
	}

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
	Case            common.UpdateCustomMetricCase
}

func (s UpdateCustomMetricStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CharacterKey")
	}

	customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CustomMetricKey")
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
	}

	switch s.Case {
	case common.CreateProperties:
		variables["properties"] = []interface{}{
			map[string]interface{}{
				"name":  "property 1",
				"type":  "NUMBER",
				"value": "20",
				"unit":  "Steps",
			},
		}

	case common.UpdateProperties:
		variables["properties"] = []interface{}{
			map[string]interface{}{
				"id":    gjson.Get(customMetric, "properties.0.id").Value(),
				"name":  "updated property 1",
				"type":  "NUMBER",
				"value": "30",
				"unit":  "CM",
			},
			map[string]interface{}{
				"name":  "new property",
				"type":  "STRING",
				"value": "ABCD",
				"unit":  "Characters",
			},
		}

	case common.DeleteProperties:
		// Remove the first property
		variables["properties"] = []interface{}{
			map[string]interface{}{
				"id":    gjson.Get(customMetric, "properties.1.id").Value(),
				"name":  "new property",
				"type":  "STRING",
				"value": "ABCD",
				"unit":  "Characters",
			},
		}
	}

	assertion := func(res *http.Response, req *http.Request) error {
		testingT, ok := (*ctx).Value(common.TestingT).(apitest.TestingT)
		if !ok {
			return common.ErrNotFoundInContext("TestingT")
		}

		json := common.ReadResponseJson(res)
		customMetric := gjson.Get(json, "data.updateCustomMetric").Value().(map[string]interface{})

		passed := assert.Empty(testingT, gjson.Get(json, "errors").Value()) &&
			assert.Equal(testingT, customMetric["name"], variables["name"]) &&
			assert.Equal(testingT, customMetric["description"], variables["description"]) &&
			assert.Equal(testingT, customMetric["style"], variables["style"]) &&
			assert.Equal(testingT, customMetric["time"], float64(0))

		properties := gjson.Get(json, "data.updateCustomMetric.properties").Value().([]interface{})
		varProperties := variables["properties"].([]interface{})

		switch s.Case {
		case common.CreateProperties:
			passed = assert.Equal(testingT, len(properties), 1) &&
				assert.NotEmpty(testingT, properties[0].(map[string]interface{})["id"])
			delete(properties[0].(map[string]interface{}), "id")
			passed = passed && assert.Equal(testingT, properties[0], varProperties[0])

		case common.UpdateProperties:
			passed = assert.Equal(testingT, len(properties), 2) &&
				assert.Equal(testingT, properties[0], varProperties[0]) &&
				assert.NotEmpty(testingT, properties[1].(map[string]interface{})["id"])
			delete(properties[1].(map[string]interface{}), "id")
			passed = passed && assert.Equal(testingT, properties[1], varProperties[1])

		case common.DeleteProperties:
			passed = assert.Equal(testingT, len(properties), 1) &&
				assert.Equal(testingT, properties[0], varProperties[0])
		}

		if !passed {
			return fmt.Errorf("assertion failed")
		}

		return nil
	}

	if s.ExpectError {
		assertion = common.AssertionError.End()
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UpdateCustomMetricQuery,
			Variables: variables,
			Assertion: assertion,
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
		return common.ErrNotFoundInContext("CharacterKey")
	}

	customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CustomMetricKey")
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
		return common.ErrNotFoundInContext("CharacterKey")
	}

	customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CustomMetricKey")
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
