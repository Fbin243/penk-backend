package mongorepo_test

// import (
// 	"testing"
// 	"time"

// 	"tenkhours/pkg/db/base"
// 	"tenkhours/services/core/entity"

// 	"github.com/stretchr/testify/assert"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// var goal = &entity.Goal{
// 	BaseEntity: &base.BaseEntity{
// 		ID:        primitive.NewObjectID().Hex(),
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	},
// 	CharacterID: primitive.NewObjectID().Hex(),
// 	Name:        "example goal",
// 	Description: "example description",
// 	StartDate:   time.Now(),
// 	EndDate:     time.Now(),
// 	Status:      entity.GoalFinishStatusUnfinished,
// 	Target:      []entity.CustomMetric{},
// }

// func TestCreateNewGoal(t *testing.T) {
// 	createdGoal, err := goalRepo.InsertOne(goal)
// 	defer cleanUpGoal(createdGoal)
// 	assert.Nil(t, err)
// 	assert.Equal(t, *createdGoal, goal)
// }

// func TestGetGoalsByCharacterID(t *testing.T) {
// 	for i := 0; i < 3; i++ {
// 		// Change the ID to create unique goals
// 		goal.ID = primitive.NewObjectID().Hex()

// 		createdGoal, err := goalRepo.InsertOne(goal)
// 		defer cleanUpGoal(createdGoal)
// 		assert.Nil(t, err)
// 	}

// 	goals, err := goalRepo.GetGoalsByCharacterID(goal.CharacterID, nil)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(goals), 3)
// 	for _, g := range goals {
// 		assert.Equal(t, g.CharacterID, goal.CharacterID)
// 	}
// }

// func cleanUpGoal(goal *entity.Goal) {
// 	goalRepo.DeleteByID(goal.ID)
// }
