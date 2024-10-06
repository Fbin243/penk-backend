package timetrackings

import (
	"log"
	"time"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
)

type CreateTimeTrackingStage struct {
	common.Metadata
	CharacterKey    common.ContextKey
	CustomMetricKey common.ContextKey
	TrackWithMetric bool
}

func (s CreateTimeTrackingStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CharacterKey")
	}

	variables := map[string]interface{}{
		"characterID": gjson.Get(character, "id").Value(),
		"startTime":   time.Now().Format(time.RFC3339Nano),
	}

	assertion := jsonpath.Chain().
		NotEqual("$.data.createTimeTracking.id", nil).
		NotEqual("$.data.createTimeTracking.characterID", nil).
		NotEqual("$.data.createTimeTracking.startTime", nil).
		Equal("$.data.createTimeTracking.endTime", nil)

	if s.TrackWithMetric {
		customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CustomMetricKey")
		}

		variables["customMetricID"] = gjson.Get(customMetric, "id").Value()
		assertion = assertion.Equal("$.data.createTimeTracking.customMetricID", variables["customMetricID"])
	} else {
		assertion = assertion.NotPresent("$.data.createTimeTracking.customMetricID")
	}

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     CreateTimeTrackingQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type UpdateTimeTracking struct {
	common.Metadata
	TimeTrackingKey common.ContextKey
	CustomMetricKey common.ContextKey
	TrackWithMetric bool
}

func (s UpdateTimeTracking) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	timeTracking, ok := (*ctx).Value(s.TimeTrackingKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("TimeTrackingKey")
	}

	variables := map[string]interface{}{
		"id": gjson.Get(timeTracking, "id").Value(),
	}

	assertion := jsonpath.Chain().
		NotEqual("$.data.updateTimeTracking.id", nil).
		NotEqual("$.data.updateTimeTracking.characterID", nil).
		NotEqual("$.data.updateTimeTracking.startTime", nil).
		NotEqual("$.data.updateTimeTracking.endTime", nil)

	if s.TrackWithMetric {
		customMetric, ok := (*ctx).Value(s.CustomMetricKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("CustomMetricKey")
		}

		variables["customMetricID"] = gjson.Get(customMetric, "id").Value()
		assertion = assertion.Equal("$.data.updateTimeTracking.customMetricID", variables["customMetricID"])
	} else {
		assertion = assertion.NotPresent("$.data.updateTimeTracking.customMetricID")
	}

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UpdateTimeTrackingQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}
