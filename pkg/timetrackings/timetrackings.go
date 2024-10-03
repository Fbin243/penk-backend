package timetrackings

import (
	"context"
	"fmt"
	"log"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/db/timetrackingsdb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimeTrackingsHandler struct {
	TimeTrackingsRepo *timetrackingsdb.TimeTrackingsRepo
	CharactersRepo    *coredb.CharactersRepo
}

func NewTimeTrackingsHandler(timeTrackingsRepo *timetrackingsdb.TimeTrackingsRepo, charactersRepo *coredb.CharactersRepo) *TimeTrackingsHandler {
	return &TimeTrackingsHandler{
		TimeTrackingsRepo: timeTrackingsRepo,
		CharactersRepo:    charactersRepo,
	}
}

func (r *TimeTrackingsHandler) GetCurrentTimeTracking(ctx context.Context, characterID primitive.ObjectID) (*timetrackingsdb.TimeTracking, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, err
	}

	if profile.ID != character.ProfileID {
		return nil, auth.ErrorPermissionDenied
	}

	result, err := r.TimeTrackingsRepo.GetCurrentTimeTrackingByCharacterID(characterID)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TimeTrackingsHandler) CreateTimeTracking(ctx context.Context, characterID primitive.ObjectID, metricID *primitive.ObjectID, startTime time.Time) (*timetrackingsdb.TimeTracking, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	serverStartTime := time.Now()
	duration := serverStartTime.Sub(startTime)
	seconds := duration.Seconds()

	// Check timeout if delay of client and server is 20 second
	if seconds > 20 {
		return nil, fmt.Errorf("server timeout, failed to start a new session")
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	if metricID != nil {
		found := false
		for _, customMetric := range character.CustomMetrics {
			if customMetric.ID == *metricID {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("custom metric does not belong to the character")
		}
	}

	// Check if the time tracking is already started
	timeTrackings, err := r.TimeTrackingsRepo.GetTimeTrackingsByCharacterID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get time trackings: %v", err)
	}

	for _, timeTracking := range timeTrackings {
		if timeTracking.EndTime.IsZero() {
			return nil, fmt.Errorf("focused session is already started")
		}
	}

	timeTracking := timetrackingsdb.TimeTracking{
		ID:              primitive.NewObjectID(),
		CharacterID:     characterID,
		CustomMetricID:  *metricID,
		StartTime:       startTime,
		MinDurationTime: 600,
		MaxDurationTime: 14400,
	}

	createdTimeTracking, err := r.TimeTrackingsRepo.CreateTimeTracking(&timeTracking)
	if err != nil {
		return nil, fmt.Errorf("failed to create time tracking: %v", err)
	}

	return createdTimeTracking, nil
}

func (r *TimeTrackingsHandler) UpdateTimeTracking(ctx context.Context, id primitive.ObjectID) (*timetrackingsdb.TimeTracking, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	endTime := time.Now()

	timeTracking, err := r.TimeTrackingsRepo.GetTimeTrackingByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get time tracking: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(timeTracking.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	if !timeTracking.EndTime.IsZero() {
		return nil, fmt.Errorf("focused session is already ended")
	}

	duration := int32(endTime.Sub(timeTracking.StartTime).Seconds())

	// TODO: TESTING
	// duration = 599 // Test for the min duration time
	// duration = 600 // Test for the min duration time
	// duration = 601 // Test for the min duration time
	// duration = 14400 // Test for the max duration time
	// duration = 14401 // Test for the max duration time

	if duration < timeTracking.MinDurationTime {
		duration = 0
		_, err := r.TimeTrackingsRepo.DeleteTimeTracking(id)
		if err != nil {
			return nil, fmt.Errorf("failed to delete time tracking: %v", err)
		}

		log.Printf("the period time is less than 10 min, so the time tracking will be deleted")
		return timeTracking, nil
	}

	if duration > timeTracking.MaxDurationTime {
		duration = int32(timeTracking.MaxDurationTime)
		log.Printf("the period time is more than 4 hours, so the time tracking will be limited to 4 hours")
	}

	timeTracking.EndTime = timeTracking.StartTime.Add(time.Duration(duration) * time.Second)

	character.TotalFocusedTime += duration
	if !timeTracking.CustomMetricID.IsZero() {
		for i, customMetric := range character.CustomMetrics {
			if customMetric.ID == timeTracking.CustomMetricID {
				character.CustomMetrics[i].Time += int32(duration)
				break
			}
		}
	}

	_, err = r.TimeTrackingsRepo.UpdateTimeTracking(timeTracking)
	if err != nil {
		return nil, fmt.Errorf("failed to update time tracking: %v", err)
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return timeTracking, nil
}
