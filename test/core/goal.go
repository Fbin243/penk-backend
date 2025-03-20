package core

import (
	"context"
	"log"

	"tenkhours/services/core/entity"
	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/tidwall/gjson"
)

type UpsertGoal struct {
	common.Metadata
	common.Case
	CharacterKey common.ContextKey
	GoalKey      common.ContextKey
}

func (s UpsertGoal) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	assertion := jsonpath.Chain().NotPresent("$.errors")
	assertions := []common.Assertion{}
	query := UpsertGoalQuery
	variables := map[string]interface{}{}

	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CharacterKey")
	}

	switch s.Case {
	case common.CreateGoal:
		metricInput1 := map[string]interface{}{
			"id":          gjson.Get(character, "metrics.0.id").Value(),
			"condition":   "gt",
			"targetValue": 123.456,
		}
		metricInput2 := map[string]interface{}{
			"id":        gjson.Get(character, "metrics.1.id").Value(),
			"condition": "ir",
			"rangeValue": map[string]interface{}{
				"min": 10.0,
				"max": 50.0,
			},
		}
		metricInput3 := map[string]interface{}{
			"id":          gjson.Get(character, "metrics.2.id").Value(),
			"condition":   "lt",
			"targetValue": 100.0,
		}
		checkbox1 := map[string]interface{}{
			"name":  "Checkbox name 1",
			"value": true,
		}
		checkbox2 := map[string]interface{}{
			"name":  "Checkbox name 2",
			"value": false,
		}
		category := map[string]interface{}{
			"id": gjson.Get(character, "categories.0.id").Value(),
			"metrics": []map[string]interface{}{
				metricInput3,
			},
		}
		goalInput := map[string]interface{}{
			"characterID": gjson.Get(character, "id").Value(),
			"name":        "Goal name",
			"description": "Goal description",
			"categories": []map[string]interface{}{
				category,
			},
			"metrics": []map[string]interface{}{
				metricInput1, metricInput2,
			},
			"checkboxes": []map[string]interface{}{
				checkbox1, checkbox2,
			},
			"startTime": "2021-01-01T00:00:00Z",
			"endTime":   "2059-12-31T23:59:59Z",
		}
		variables["input"] = goalInput
		assertion = assertion.
			NotEqual("$.data.upsertGoal.id", "").
			Equal("$.data.upsertGoal.characterID", goalInput["characterID"]).
			Equal("$.data.upsertGoal.name", goalInput["name"]).
			Equal("$.data.upsertGoal.description", goalInput["description"]).
			Equal("$.data.upsertGoal.startTime", goalInput["startTime"]).
			Equal("$.data.upsertGoal.endTime", goalInput["endTime"]).
			Equal("$.data.upsertGoal.categories[0].id", category["id"]).
			Equal("$.data.upsertGoal.categories[0].metrics[0].id", metricInput3["id"]).
			Equal("$.data.upsertGoal.categories[0].metrics[0].condition", metricInput3["condition"]).
			Equal("$.data.upsertGoal.categories[0].metrics[0].targetValue", metricInput3["targetValue"]).
			Equal("$.data.upsertGoal.metrics[0].id", metricInput1["id"]).
			Equal("$.data.upsertGoal.metrics[0].condition", metricInput1["condition"]).
			Equal("$.data.upsertGoal.metrics[0].targetValue", metricInput1["targetValue"]).
			Equal("$.data.upsertGoal.metrics[1].id", metricInput2["id"]).
			Equal("$.data.upsertGoal.metrics[1].condition", metricInput2["condition"]).
			Equal("$.data.upsertGoal.metrics[1].rangeValue", metricInput2["rangeValue"]).
			NotEqual("$.data.upsertGoal.checkboxes[0].id", nil).
			Equal("$.data.upsertGoal.checkboxes[0].name", checkbox1["name"]).
			Equal("$.data.upsertGoal.checkboxes[0].value", checkbox1["value"]).
			NotEqual("$.data.upsertGoal.checkboxes[1].id", nil).
			Equal("$.data.upsertGoal.checkboxes[1].name", checkbox2["name"]).
			Equal("$.data.upsertGoal.checkboxes[1].value", checkbox2["value"]).
			Equal("$.data.upsertGoal.status", string(entity.GoalStatusPlanned))
		assertions = append(assertions,
			jsonpath.Len("$.data.upsertGoal.categories", 1),
			jsonpath.Len("$.data.upsertGoal.categories[0].metrics", 1),
			jsonpath.Len("$.data.upsertGoal.metrics", 2),
			jsonpath.Len("$.data.upsertGoal.checkboxes", 2))

	case common.UpdateGoal:
		goal, ok := (*ctx).Value(s.GoalKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("GoalKey")
		}
		metricInput1 := map[string]interface{}{
			"id":          gjson.Get(character, "metrics.0.id").Value(),
			"condition":   "lte",
			"targetValue": 3.0,
		}
		metricInput2 := map[string]interface{}{
			"id":        gjson.Get(character, "metrics.1.id").Value(),
			"condition": "ir",
			"rangeValue": map[string]interface{}{
				"min": 1.0,
				"max": 3.0,
			},
		}
		metricInput3 := map[string]interface{}{
			"id":          gjson.Get(character, "metrics.2.id").Value(),
			"condition":   "gte",
			"targetValue": 1.0,
		}
		checkbox1 := map[string]interface{}{
			"id":    gjson.Get(goal, "checkboxes.0.id").Value(),
			"name":  "Update checkbox name 1",
			"value": true,
		}
		category := map[string]interface{}{
			"id": gjson.Get(character, "categories.0.id").Value(),
			"metrics": []map[string]interface{}{
				metricInput3,
			},
		}
		goalInput := map[string]interface{}{
			"id":          gjson.Get(goal, "id").Value(),
			"characterID": gjson.Get(character, "id").Value(),
			"name":        "Update goal name",
			"description": "Update goal description",
			"categories": []map[string]interface{}{
				category,
			},
			"metrics": []map[string]interface{}{
				metricInput1, metricInput2,
			},
			"checkboxes": []map[string]interface{}{
				checkbox1,
			},
			"startTime": "2025-02-02T00:00:00Z",
			"endTime":   "2069-09-16T23:59:59Z",
		}
		variables["input"] = goalInput
		assertion = assertion.
			NotEqual("$.data.upsertGoal.id", "").
			Equal("$.data.upsertGoal.characterID", goalInput["characterID"]).
			Equal("$.data.upsertGoal.name", goalInput["name"]).
			Equal("$.data.upsertGoal.description", goalInput["description"]).
			Equal("$.data.upsertGoal.startTime", goalInput["startTime"]).
			Equal("$.data.upsertGoal.endTime", goalInput["endTime"]).
			Equal("$.data.upsertGoal.categories[0].id", category["id"]).
			Equal("$.data.upsertGoal.categories[0].metrics[0].id", metricInput3["id"]).
			Equal("$.data.upsertGoal.categories[0].metrics[0].condition", metricInput3["condition"]).
			Equal("$.data.upsertGoal.categories[0].metrics[0].targetValue", metricInput3["targetValue"]).
			Equal("$.data.upsertGoal.metrics[0].id", metricInput1["id"]).
			Equal("$.data.upsertGoal.metrics[0].condition", metricInput1["condition"]).
			Equal("$.data.upsertGoal.metrics[0].targetValue", metricInput1["targetValue"]).
			Equal("$.data.upsertGoal.metrics[1].id", metricInput2["id"]).
			Equal("$.data.upsertGoal.metrics[1].condition", metricInput2["condition"]).
			Equal("$.data.upsertGoal.metrics[1].rangeValue", metricInput2["rangeValue"]).
			Present("$.data.upsertGoal.checkboxes[0].id").
			Equal("$.data.upsertGoal.checkboxes[0].name", checkbox1["name"]).
			Equal("$.data.upsertGoal.checkboxes[0].value", checkbox1["value"]).
			Equal("$.data.upsertGoal.status", string(entity.GoalStatusCompleted))
		assertions = append(assertions,
			jsonpath.Len("$.data.upsertGoal.categories", 1),
			jsonpath.Len("$.data.upsertGoal.categories[0].metrics", 1),
			jsonpath.Len("$.data.upsertGoal.metrics", 2),
			jsonpath.Len("$.data.upsertGoal.checkboxes", 1))
	}

	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     query,
			Variables: variables,
			Assertion: append(assertions, assertion.End()),
		})
}

type DeleteGoal struct {
	common.Metadata
	GoalKey common.ContextKey
}

func (s DeleteGoal) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	goal, ok := (*ctx).Value(s.GoalKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("GoalKey")
	}

	variables := map[string]interface{}{
		"id": gjson.Get(goal, "id").Value(),
	}

	assertion := jsonpath.Chain().NotPresent("$.errors")
	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     DeleteGoalQuery,
			Variables: variables,
			Assertion: []common.Assertion{assertion.End()},
		})
}
