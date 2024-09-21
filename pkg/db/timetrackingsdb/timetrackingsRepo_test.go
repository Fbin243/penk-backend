package timetrackingsdb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type timeTrackingInputType struct {
	CharacterID     primitive.ObjectID
	CustomMetricID  primitive.ObjectID
	StartTime       time.Time
	EndTime         time.Time
	MinDurationTime int32
	MaxDurationTime int32
}

var timeTrackingInput = &timeTrackingInputType{
	CharacterID:     primitive.NewObjectID(),
	CustomMetricID:  primitive.NewObjectID(),
	StartTime:       time.Now().UTC().Add(-1 * time.Hour),
	EndTime:         time.Now().UTC(),
	MinDurationTime: 30,
	MaxDurationTime: 120,
}

func newTimeTrackingFromInput(input *timeTrackingInputType) *TimeTracking {
	return &TimeTracking{
		ID:              primitive.NewObjectID(),
		CharacterID:     input.CharacterID,
		CustomMetricID:  input.CustomMetricID,
		StartTime:       input.StartTime,
		EndTime:         input.EndTime,
		MinDurationTime: input.MinDurationTime,
		MaxDurationTime: input.MaxDurationTime,
	}
}

func assertWithTimeTrackingInput(t *testing.T, timeTracking *TimeTracking, input *timeTrackingInputType) {
	assert.Equal(t, timeTracking.CharacterID, input.CharacterID)
	assert.Equal(t, timeTracking.CustomMetricID, input.CustomMetricID)
	assert.WithinDuration(t, timeTracking.StartTime, input.StartTime, 1*time.Second)
	assert.WithinDuration(t, timeTracking.EndTime, input.EndTime, 1*time.Second)
	assert.Equal(t, timeTracking.MinDurationTime, input.MinDurationTime)
	assert.Equal(t, timeTracking.MaxDurationTime, input.MaxDurationTime)
}

func setupTest(t *testing.T) (*TimeTracking, func()) {
	_, err := timeTrackingsRepo.Collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		t.Fatalf("Failed to clean up collection: %v", err)
	}

	timeTracking := newTimeTrackingFromInput(timeTrackingInput)

	_, err = timeTrackingsRepo.CreateTimeTracking(timeTracking)
	if err != nil {
		t.Fatalf("Failed to create time tracking: %v", err)
	}

	cleanup := func() {
		_, err := timeTrackingsRepo.DeleteTimeTracking(timeTracking.ID)
		if err != nil {
			t.Fatalf("Failed to delete time tracking: %v", err)
		}
	}

	return timeTracking, cleanup
}

func TestCreateTimeTracking(t *testing.T) {
	timeTracking := newTimeTrackingFromInput(timeTrackingInput)

	createdTimeTracking, err := timeTrackingsRepo.CreateTimeTracking(timeTracking)
	assert.Nil(t, err)
	assertWithTimeTrackingInput(t, createdTimeTracking, timeTrackingInput)
}

func TestGetTimeTrackingByID(t *testing.T) {
	timeTracking, cleanup := setupTest(t)
	defer cleanup()

	queriedTimeTracking, err := timeTrackingsRepo.GetTimeTrackingByID(timeTracking.ID)
	assert.Nil(t, err)
	assertWithTimeTrackingInput(t, queriedTimeTracking, timeTrackingInput)
}

func TestGetTimeTrackingsByCharacterID(t *testing.T) {
	timeTracking, cleanup := setupTest(t)
	defer cleanup()

	timeTrackings, err := timeTrackingsRepo.GetTimeTrackingsByCharacterID(timeTracking.CharacterID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(timeTrackings))
	assertWithTimeTrackingInput(t, &timeTrackings[0], timeTrackingInput)
}

// func TestUpdateTimeTracking(t *testing.T) {
// 	timeTracking, cleanup := setupTest(t)
// 	defer cleanup()

// 	fixedStartTime := time.Now().UTC().Add(-2 * time.Hour)
// 	fixedEndTime := time.Now().UTC()

// 	updateInput := &timeTrackingInputType{
// 		CharacterID:     timeTracking.CharacterID,
// 		CustomMetricID:  timeTracking.CustomMetricID,
// 		StartTime:       fixedStartTime,
// 		EndTime:         fixedEndTime,
// 		MinDurationTime: 60,
// 		MaxDurationTime: 180,
// 	}

// 	timeTracking.StartTime = updateInput.StartTime
// 	timeTracking.EndTime = updateInput.EndTime
// 	timeTracking.MinDurationTime = updateInput.MinDurationTime
// 	timeTracking.MaxDurationTime = updateInput.MaxDurationTime

// 	_, err := timeTrackingsRepo.UpdateTimeTracking(timeTracking)
// 	assert.Nil(t, err)

// 	updatedTimeTracking, err := timeTrackingsRepo.GetTimeTrackingByID(timeTracking.ID)
// 	assert.Nil(t, err)

// 	assert.WithinDuration(t, updatedTimeTracking.StartTime, updateInput.StartTime, 1*time.Second)
// 	assert.WithinDuration(t, updatedTimeTracking.EndTime, updateInput.EndTime, 1*time.Second)
// 	assert.Equal(t, updateInput.MinDurationTime, updatedTimeTracking.MinDurationTime)
// 	assert.Equal(t, updateInput.MaxDurationTime, updatedTimeTracking.MaxDurationTime)
// }

func TestDeleteTimeTracking(t *testing.T) {
	timeTracking, cleanup := setupTest(t)
	defer cleanup()

	deletedTimeTracking, err := timeTrackingsRepo.DeleteTimeTracking(timeTracking.ID)
	assert.Nil(t, err)
	assertWithTimeTrackingInput(t, deletedTimeTracking, timeTrackingInput)

	queriedTimeTracking, err := timeTrackingsRepo.GetTimeTrackingByID(timeTracking.ID)
	assert.Nil(t, err)
	assert.Nil(t, queriedTimeTracking)
}
