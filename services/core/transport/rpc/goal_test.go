package rpc_test

import (
	"testing"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
	"tenkhours/services/core/transport/rpc"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var goalCheckbox = entity.Checkbox{
	ID:    mongodb.GenObjectID(),
	Name:  "Goal name",
	Value: true,
}

var goalMetric = entity.GoalMetric{
	ID:          mongodb.GenObjectID(),
	Condition:   entity.MetricConditionInRange,
	TargetValue: lo.ToPtr(100.0),
	RangeValue: &entity.Range{
		Min: 200,
		Max: 300,
	},
}

var goal = entity.Goal{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	},
	CharacterID: mongodb.GenObjectID(),
	Name:        "Goal name",
	Description: "Goal desc",
	StartTime:   utils.Now(),
	EndTime:     utils.Now(),
	Metrics: []entity.GoalMetric{
		goalMetric, goalMetric,
	},
	Checkboxes: []entity.Checkbox{
		goalCheckbox, goalCheckbox,
	},
}

func TestMapGoal(t *testing.T) {
	rpcGoal, err := rpc.Map[entity.Goal, core.Goal](&goal, append(rpc.UnixTimeConverter, rpc.MetricConditionConverter...))
	assert.NoError(t, err)

	assert.Equal(t, goal.ID, rpcGoal.Id)
	assert.Equal(t, goal.CharacterID, rpcGoal.CharacterId)
	assert.Equal(t, goal.Name, rpcGoal.Name)
	assert.Equal(t, goal.Description, rpcGoal.Description)
	assert.Equal(t, goal.StartTime.Unix(), rpcGoal.StartTime)
	assert.Equal(t, goal.EndTime.Unix(), rpcGoal.EndTime)
	assert.Equal(t, goal.CreatedAt.Unix(), rpcGoal.CreatedAt)
	assert.Equal(t, goal.UpdatedAt.Unix(), rpcGoal.UpdatedAt)
	for i, expectedMetric := range goal.Metrics {
		assertMetric(t, expectedMetric, rpcGoal.Metrics[i])
	}
	for i, expectedCheckbox := range goal.Checkboxes {
		assertCheckbox(t, expectedCheckbox, rpcGoal.Checkboxes[i])
	}
}

func assertMetric(t *testing.T, expected entity.GoalMetric, actual *core.GoalMetric) {
	assert.Equal(t, expected.ID, actual.Id)
	assert.Equal(t, string(expected.Condition), actual.Condition.String())
	assert.Equal(t, *expected.TargetValue, *actual.TargetValue)
	assert.Equal(t, expected.RangeValue.Min, actual.RangeValue.Min)
	assert.Equal(t, expected.RangeValue.Max, actual.RangeValue.Max)
}

func assertCheckbox(t *testing.T, expected entity.Checkbox, actual *core.Checkbox) {
	assert.Equal(t, expected.ID, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Value, actual.Value)
}

// Input
var rpcGoalMetricInput = &core.GoalMetricInput{
	Id:          mongodb.GenObjectID(),
	Condition:   core.MetricCondition_ir,
	TargetValue: lo.ToPtr(100.0),
	RangeValue: &core.RangeInput{
		Min: 200,
		Max: 300,
	},
}

var rpcGoalCheckboxInput = &core.CheckboxInput{
	Id:    lo.ToPtr(mongodb.GenObjectID()),
	Name:  "Checkbox name",
	Value: true,
}

var rpcGoalInput = &core.GoalInput{
	Id:          lo.ToPtr(mongodb.GenObjectID()),
	Name:        "Goal name",
	Description: lo.ToPtr("Goal desc"),
	StartTime:   utils.Now().Unix(),
	EndTime:     utils.Now().Unix(),
	Metrics: []*core.GoalMetricInput{
		rpcGoalMetricInput, rpcGoalMetricInput,
	},
	Checkboxes: []*core.CheckboxInput{
		rpcGoalCheckboxInput, rpcGoalCheckboxInput,
	},
}

func TestMapGoalInput(t *testing.T) {
	goalInput, err := rpc.Map[core.GoalInput, entity.GoalInput](rpcGoalInput, append(rpc.UnixTimeConverter, rpc.MetricConditionConverter...))
	assert.NoError(t, err)

	assert.Equal(t, rpcGoalInput.Id, goalInput.ID)
	assert.Equal(t, rpcGoalInput.Name, goalInput.Name)
	assert.Equal(t, *rpcGoalInput.Description, *goalInput.Description)
	assert.Equal(t, rpcGoalInput.StartTime, goalInput.StartTime.Unix())
	assert.Equal(t, rpcGoalInput.EndTime, goalInput.EndTime.Unix())

	for i, expectedMetric := range rpcGoalInput.Metrics {
		assertMetricInput(t, expectedMetric, goalInput.Metrics[i])
	}
	for i, expectedCheckbox := range rpcGoalInput.Checkboxes {
		assertCheckboxInput(t, expectedCheckbox, goalInput.Checkboxes[i])
	}
}

func assertMetricInput(t *testing.T, expected *core.GoalMetricInput, actual entity.GoalMetricInput) {
	assert.Equal(t, expected.Id, actual.ID)
	assert.Equal(t, expected.Condition.String(), string(actual.Condition))
	assert.Equal(t, *expected.TargetValue, *actual.TargetValue)
	assert.Equal(t, expected.RangeValue.Min, actual.RangeValue.Min)
	assert.Equal(t, expected.RangeValue.Max, actual.RangeValue.Max)
}

func assertCheckboxInput(t *testing.T, expected *core.CheckboxInput, actual entity.CheckboxInput) {
	assert.Equal(t, expected.Id, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Value, actual.Value)
}
