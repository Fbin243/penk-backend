package repo_test

import (
	"testing"
	"time"

	"tenkhours/pkg/db"
	"tenkhours/services/core/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var goal = &repo.Goal{
	BaseModel: &db.BaseModel{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	CharacterID: primitive.NewObjectID(),
	Name:        "example goal",
	Description: "example description",
	StartDate:   time.Now(),
	EndDate:     time.Now(),
	Status:      repo.GoalStatusActive,
	Target:      []repo.CustomMetric{},
}

func TestCreateNewGoal(t *testing.T) {
	createdGoal, err := goalsRepo.InsertOne(goal)
	defer cleanUpGoal(createdGoal)
	assert.Nil(t, err)
	assert.Equal(t, *createdGoal, goal)
}

func TestGetGoalsByCharacterID(t *testing.T) {
	for i := 0; i < 3; i++ {
		// Change the ID to create unique goals
		goal.BaseModel.ID = primitive.NewObjectID()

		createdGoal, err := goalsRepo.InsertOne(goal)
		defer cleanUpGoal(createdGoal)
		assert.Nil(t, err)
	}

	goals, err := goalsRepo.GetGoalsByCharacterID(goal.CharacterID)
	assert.Nil(t, err)
	assert.Equal(t, len(goals), 3)
	for _, g := range goals {
		assert.Equal(t, g.CharacterID, goal.CharacterID)
	}
}

func cleanUpGoal(goal *repo.Goal) {
	goalsRepo.DeleteById(goal.ID)
}
