package rpc_test

import (
	"fmt"
	"testing"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var goalCheckbox = entity.Checkbox{
	ID:    mongodb.GenObjectID(),
	Name:  "Goal name",
	Value: true,
}

var goalMetric = entity.GoalMetric{
	Metric: &entity.Metric{
		ID:    mongodb.GenObjectID(),
		Name:  "Metric name",
		Value: 100,
		Unit:  "Metric unit",
	},
	Condition:   entity.MetricConditionEqual,
	TargetValue: lo.ToPtr(100.0),
	RangeValue: &entity.Range{
		Min: 200,
		Max: 300,
	},
}

var goalCategory = entity.GoalCategory{
	Category: &entity.Category{
		ID:          mongodb.GenObjectID(),
		Name:        "Category name",
		Description: "Category desc",
		Style: entity.CategoryStyle{
			Color: "#000000",
			Icon:  "icon.png",
		},
	},
	Metrics: []entity.GoalMetric{
		goalMetric, goalMetric,
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
	Status:      entity.GoalFinishStatusFinished,
	Target: entity.GoalTarget{
		Categories: []entity.GoalCategory{
			goalCategory, goalCategory,
		},
		Metrics: []entity.GoalMetric{
			goalMetric, goalMetric,
		},
		Checkboxes: []entity.Checkbox{
			goalCheckbox, goalCheckbox,
		},
	},
}

func TestMapGoal(t *testing.T) {
	rpcGoal := &core.Goal{}

	copier.Copy(rpcGoal, &goal)
	rpcGoal.Status = core.GoalFinishStatus(core.GoalFinishStatus_value[string(goal.Status)])
	rpcGoal.CreatedAt = goal.CreatedAt.Unix()
	rpcGoal.UpdatedAt = goal.UpdatedAt.Unix()
	rpcGoal.StartTime = goal.StartTime.Unix()
	rpcGoal.EndTime = goal.EndTime.Unix()
	copier.Copy(&rpcGoal.Categories, &goal.Target.Categories)
	copier.Copy(&rpcGoal.Metrics, &goal.Target.Metrics)
	copier.Copy(&rpcGoal.Checkboxes, &goal.Target.Checkboxes)

	assert.Equal(t, goal.ID, rpcGoal.Id)
	assert.Equal(t, goal.CharacterID, rpcGoal.CharacterId)
	assert.Equal(t, goal.Name, rpcGoal.Name)
	assert.Equal(t, goal.Description, rpcGoal.Description)
	assert.Equal(t, string(goal.Status), rpcGoal.Status.String())
	assert.Equal(t, goal.StartTime.Unix(), rpcGoal.StartTime)
	assert.Equal(t, goal.EndTime.Unix(), rpcGoal.EndTime)
	assert.Equal(t, goal.CreatedAt.Unix(), rpcGoal.CreatedAt)
	assert.Equal(t, goal.UpdatedAt.Unix(), rpcGoal.UpdatedAt)
	for i, expectedCategory := range goal.Target.Categories {
		assertCategory(t, expectedCategory, rpcGoal.Categories[i])
	}
	for i, expectedMetric := range goal.Target.Metrics {
		assertMetric(t, expectedMetric, rpcGoal.Metrics[i])
	}
	for i, expectedCheckbox := range goal.Target.Checkboxes {
		assertCheckbox(t, expectedCheckbox, rpcGoal.Checkboxes[i])
	}
}

func assertCategory(t *testing.T, expected entity.GoalCategory, actual *core.GoalCategory) {
	assert.Equal(t, expected.ID, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.Style.Color, actual.Style.Color)
	assert.Equal(t, expected.Style.Icon, actual.Style.Icon)
	for i, expectedMetric := range expected.Metrics {
		assertMetric(t, expectedMetric, actual.Metrics[i])
	}
}

func assertMetric(t *testing.T, expected entity.GoalMetric, actual *core.GoalMetric) {
	assert.Equal(t, expected.ID, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Value, float64(actual.Value))
	assert.Equal(t, expected.Unit, actual.Unit)
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

var rpcGoalCategoryInput = &core.GoalCategoryInput{
	Id: mongodb.GenObjectID(),
	Metrics: []*core.GoalMetricInput{
		rpcGoalMetricInput, rpcGoalMetricInput,
	},
}

var rpcGoalInput = &core.GoalInput{
	Id:          lo.ToPtr(mongodb.GenObjectID()),
	CharacterId: mongodb.GenObjectID(),
	Name:        "Goal name",
	Description: lo.ToPtr("Goal desc"),
	StartTime:   utils.Now().Unix(),
	EndTime:     utils.Now().Unix(),
	Categories: []*core.GoalCategoryInput{
		rpcGoalCategoryInput, rpcGoalCategoryInput,
	},
	Metrics: []*core.GoalMetricInput{
		rpcGoalMetricInput, rpcGoalMetricInput,
	},
	Checkboxes: []*core.CheckboxInput{
		rpcGoalCheckboxInput, rpcGoalCheckboxInput,
	},
}

func TestMapGoalInput(t *testing.T) {
	goalInput := &entity.GoalInput{}

	copier.CopyWithOption(goalInput, rpcGoalInput, copier.Option{
		IgnoreEmpty: true,
		Converters: []copier.TypeConverter{
			{
				SrcType: core.MetricCondition(0),
				DstType: entity.MetricCondition(fmt.Sprint(0)),
				Fn: func(src any) (any, error) {
					return entity.MetricCondition(
						core.MetricCondition_name[int32(src.(core.MetricCondition).Number())],
					), nil
				},
			},
		},
	})
	goalInput.EndTime = utils.UnixToTime(rpcGoalInput.EndTime)
	goalInput.StartTime = utils.UnixToTime(rpcGoalInput.StartTime)

	assert.Equal(t, rpcGoalInput.Id, goalInput.ID)
	assert.Equal(t, rpcGoalInput.CharacterId, goalInput.CharacterID)
	assert.Equal(t, rpcGoalInput.Name, goalInput.Name)
	assert.Equal(t, *rpcGoalInput.Description, *goalInput.Description)
	assert.Equal(t, rpcGoalInput.StartTime, goalInput.StartTime.Unix())
	assert.Equal(t, rpcGoalInput.EndTime, goalInput.EndTime.Unix())

	for i, expectedCategory := range rpcGoalInput.Categories {
		assertCategoryInput(t, expectedCategory, goalInput.Categories[i])
	}
	for i, expectedMetric := range rpcGoalInput.Metrics {
		assertMetricInput(t, expectedMetric, goalInput.Metrics[i])
	}
	for i, expectedCheckbox := range rpcGoalInput.Checkboxes {
		assertCheckboxInput(t, expectedCheckbox, goalInput.Checkboxes[i])
	}
}

func assertCategoryInput(t *testing.T, expected *core.GoalCategoryInput, actual entity.GoalCategoryInput) {
	assert.Equal(t, expected.Id, actual.ID)
	for i, expectedMetric := range expected.Metrics {
		assertMetricInput(t, expectedMetric, actual.Metrics[i])
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
