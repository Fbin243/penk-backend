package core

import (
	"log"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
)

type CreateMetricPropertyStage struct {
	common.Metadata
	CustomMetricKey common.ContextKey
	CharacterKey    common.ContextKey
}

func (s CreateMetricPropertyStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CustomMetricKey)
	}

	metricPropertyInput := map[string]interface{}{
		"name":  "property 2",
		"type":  "NUMBER",
		"value": "value 2",
	}

	variables := map[string]interface{}{
		"characterID": gjson.Get(character, "id").Value(),
		"metricID":    gjson.Get(customMetric, "id").Value(),
		"input":       metricPropertyInput,
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Present("$.data.createMetricProperty.id").
		Equal("$.data.createMetricProperty.name", metricPropertyInput["name"]).
		Equal("$.data.createMetricProperty.type", metricPropertyInput["type"]).
		Equal("$.data.createMetricProperty.value", metricPropertyInput["value"]).
		Equal("$.data.createMetricProperty.unit", metricPropertyInput["unit"])

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     CreateMetricPropertyQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type UpdateMetricPropertyStage struct {
	common.Metadata
	MetricPropertyKey common.ContextKey
	CustomMetricKey   common.ContextKey
	CharacterKey      common.ContextKey
}

func (s UpdateMetricPropertyStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CustomMetricKey)
	}

	metricProperty, ok := (*ctx).Value(s.MetricPropertyKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.MetricPropertyKey)
	}

	metricPropertyInput := map[string]interface{}{
		"name":  "updated property 2",
		"type":  "STRING",
		"value": "updated value 2",
		"unit":  "unit",
	}

	variables := map[string]interface{}{
		"id":          gjson.Get(metricProperty, "id").Value(),
		"characterID": gjson.Get(character, "id").Value(),
		"metricID":    gjson.Get(customMetric, "id").Value(),
		"input":       metricPropertyInput,
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.updateMetricProperty.id", variables["id"]).
		Equal("$.data.updateMetricProperty.name", metricPropertyInput["name"]).
		Equal("$.data.updateMetricProperty.type", metricPropertyInput["type"]).
		Equal("$.data.updateMetricProperty.value", metricPropertyInput["value"]).
		Equal("$.data.updateMetricProperty.unit", metricPropertyInput["unit"])

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UpdateMetricPropertyQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type DeleteMetricPropertyStage struct {
	common.Metadata
	MetricPropertyKey common.ContextKey
	CustomMetricKey   common.ContextKey
	CharacterKey      common.ContextKey
}

func (s DeleteMetricPropertyStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CharacterKey)
	}

	customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.CustomMetricKey)
	}

	metricProperty, ok := (*ctx).Value(s.MetricPropertyKey).(string)
	if !ok {
		return common.ErrNotFoundInContext(s.MetricPropertyKey)
	}

	variables := map[string]interface{}{
		"id":          gjson.Get(metricProperty, "id").Value(),
		"characterID": gjson.Get(character, "id").Value(),
		"metricID":    gjson.Get(customMetric, "id").Value(),
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.deleteMetricProperty.id", variables["id"])

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     DeleteMetricPropertyQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}
