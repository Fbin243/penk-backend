package timetracking

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
	common.Case
	CharacterKey common.ContextKey
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

	assertion := jsonpath.Chain().NotPresent("$.errors").
		NotEqual("$.data.createTimeTracking.id", nil).
		NotEqual("$.data.createTimeTracking.characterID", nil).
		NotEqual("$.data.createTimeTracking.startTime", nil).
		Equal("$.data.createTimeTracking.endTime", nil)

	switch s.Case {
	case common.TimeTrackingWithoutCategory:
		assertion = assertion.Equal("$.data.createTimeTracking.categoryID", nil)

	case common.TimeTrackingWithCategory:
		// Track with the first metric
		assertion = assertion.Equal("$.data.createTimeTracking.categoryID", gjson.Get(character, "categories.0.id").Value())
	}

	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     CreateTimeTrackingQuery,
			Variables: variables,
			Assertion: []common.Assertion{assertion.End()},
		})
}

type GetCurrentTimeTrackingStage struct {
	common.Metadata
	common.Case
	CharacterKey common.ContextKey
}

func (s GetCurrentTimeTrackingStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CharacterKey")
	}

	assertion := jsonpath.Chain().NotPresent("$.errors")
	if s.Case == common.CurrentTimeTrackingExist {
		assertion.
			NotEqual("$.data.currentTimeTracking.id", nil).
			Equal("$.data.currentTimeTracking.characterID", gjson.Get(character, "id").Value()).
			Equal("$.data.currentTimeTracking.categoryID", gjson.Get(character, "categories.0.id").Value()).
			NotEqual("$.data.currentTimeTracking.startTime", nil).
			Equal("$.data.currentTimeTracking.endTime", nil)
	} else {
		assertion.Equal("$.data.currentTimeTracking", nil)
	}

	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     GetCurrentTimeTrackingQuery,
			Assertion: []common.Assertion{assertion.End()},
		})
}

type UpdateTimeTracking struct {
	common.Metadata
	common.Case
	CharacterKey    common.ContextKey
	TimeTrackingKey common.ContextKey
}

func (s UpdateTimeTracking) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	timeTracking, ok := (*ctx).Value(s.TimeTrackingKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("TimeTrackingKey")
	}

	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CharacterKey")
	}

	assertion := jsonpath.Chain().
		NotPresent("$.errors").
		Equal("$.data.updateTimeTracking.timeTracking.id", gjson.Get(timeTracking, "id").Value()).
		NotEqual("$.data.updateTimeTracking.timeTracking.characterID", nil).
		NotEqual("$.data.updateTimeTracking.timeTracking.startTime", nil).
		NotEqual("$.data.updateTimeTracking.timeTracking.endTime", nil)

	switch s.Case {
	case common.TimeTrackingWithoutCategory:
		assertion = assertion.Equal("$.data.updateTimeTracking.timeTracking.categoryID", nil).
			NotEqual("$.data.updateTimeTracking.gold", nil).
			NotEqual("$.data.updateTimeTracking.normal", nil)

	case common.TimeTrackingWithCategory:
		// Track with the first category
		assertion = assertion.
			Equal("$.data.updateTimeTracking.timeTracking.categoryID", gjson.Get(character, "categories.0.id").Value()).
			NotEqual("$.data.updateTimeTracking.gold", nil).
			NotEqual("$.data.updateTimeTracking.normal", nil)
	}

	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     UpdateTimeTrackingQuery,
			Assertion: []common.Assertion{assertion.End()},
		})
}
