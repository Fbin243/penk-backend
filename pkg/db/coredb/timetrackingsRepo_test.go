package coredb

// import (
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type timeInputType struct {
// 	StartTime time.Time
// }

// type timeUpdateType struct {
// 	EndTime time.Time
// }

// var timeInput = &timeInputType{
// 	StartTime: time.Now(),
// }

// var timeUpdateInput = &timeUpdateType{
// 	EndTime: time.Now(),
// }

// func newTimeFromInput(input *timeInputType) *TimeTracking {
// 	return &TimeTracking{
// 		ID:             primitive.NewObjectID(),
// 		CharacterID:    primitive.NewObjectID(),
// 		CustomMetricID: primitive.NewObjectID(),
// 		StartTime:      time.Now(),
// 	}
// }

// func newUpdateTimeFromInput(input *timeUpdateType) *TimeTracking {
// 	return &TimeTracking{
// 		ID:             primitive.NewObjectID(),
// 		CharacterID:    primitive.NewObjectID(),
// 		CustomMetricID: primitive.NewObjectID(),
// 		EndTime:        time.Now(),
// 	}
// }

// func assertWithTimeInput(t *testing.T, time *TimeTracking, input *timeInputType) {
// 	assert.Equal(t, time.StartTime, input.StartTime)
// }

// func assertWithTimeUpdateInput(t *testing.T, time *TimeTracking, input *timeUpdateType) {
// 	assert.Equal(t, time.StartTime, input.EndTime)
// }

// func TestCreateNewTime(t *testing.T) {
// 	time := newTimeFromInput(timeInput)

// 	createdTime, err := TimeTrackingsRepo.CreateTimeTracking(time)
// 	assert.Nil(t, err)
// 	assertWithTimeInput(t, createdTime, timeInput)
// }

// func TestUpdateTime(t *testing.T) {
// 	time := newTimeFromInput(timeInput)

// 	_, err := TimeTrackingsRepo.CreateTimeTracking(time)
// 	assert.Nil(t, err)

// 	updateInput := newUpdateTimeFromInput(timeUpdateInput)

// 	time.EndTime = updateInput.EndTime
// 	updatedTime, err := TimeTrackingsRepo.UpdateTimeTracking(time)
// 	assert.Nil(t, err)
// 	assertWithTimeUpdateInput(t, updatedTime, updateInput)
// }
