package mongorepo_test

import (
	"context"
	"testing"
	"time"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var goal = &entity.Goal{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	CharacterID: mongodb.GenObjectID(),
	Name:        "example goal",
	Description: "example description",
	StartTime:   time.Now(),
	EndTime:     time.Now(),
	Metrics:     []entity.GoalMetric{targetMetric, rangeMetric, targetMetric},
	Checkboxes:  []entity.Checkbox{checkbox, checkbox, checkbox},
}

var targetMetric = entity.GoalMetric{
	ID:          mongodb.GenObjectID(),
	Condition:   entity.MetricConditionGreaterThan,
	TargetValue: lo.ToPtr(10.0),
}

var rangeMetric = entity.GoalMetric{
	ID:         mongodb.GenObjectID(),
	Condition:  entity.MetricConditionInRange,
	RangeValue: &entity.Range{Min: 0, Max: 10},
}

var checkbox = entity.Checkbox{
	ID:    mongodb.GenObjectID(),
	Name:  "Checkbox name",
	Value: false,
}

func TestCreateNewGoal(t *testing.T) {
	createdGoal, err := goalRepo.InsertOne(context.Background(), goal)
	defer cleanUpGoal(t, createdGoal.ID)
	assert.Nil(t, err)
	assert.Equal(t, goal, createdGoal)
}

func TestGetGoalsByCharacterID(t *testing.T) {
	for range 3 {
		goal.ID = mongodb.GenObjectID()
		_, err := goalRepo.InsertOne(context.Background(), goal)
		defer cleanUpGoal(t, goal.ID)
		assert.Nil(t, err)
	}

	goals, err := goalRepo.Find(context.Background(), entity.GoalPipeline{
		Filter: &entity.GoalFilter{
			CharacterID: &goal.CharacterID,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(goals))
	for _, g := range goals {
		assert.Equal(t, g.CharacterID, goal.CharacterID)
	}
}

func cleanUpGoal(t assert.TestingT, id string) {
	_, err := goalRepo.FindAndDeleteByID(context.Background(), id)
	assert.Nil(t, err)
}
