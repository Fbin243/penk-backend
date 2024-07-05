package timetrackings

import (
	"fmt"
	"log"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeTrackingsResolver struct {
	TimeTrackingsRepo *coredb.TimeTrackingsRepo
	CharactersRepo    *coredb.CharactersRepo
}

func NewTimeTrackingsResolver(timeTrackingsRepo *coredb.TimeTrackingsRepo, charactersRepo *coredb.CharactersRepo) *TimeTrackingsResolver {
	return &TimeTrackingsResolver{
		TimeTrackingsRepo: timeTrackingsRepo,
		CharactersRepo:    charactersRepo,
	}
}

func (r *TimeTrackingsResolver) CreateTimeTracking(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.UserID != user.ID {
		return nil, fmt.Errorf("permission denied")
	}

	customMetricID, ok := params.Args["customMetricID"].(string)
	customMetricOID := primitive.ObjectID{}
	if ok {
		customMetricOID, err = primitive.ObjectIDFromHex(customMetricID)
		if err != nil {
			return nil, err
		}

		found := false
		for _, customMetric := range character.CustomMetrics {
			if customMetric.ID == customMetricOID {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("custom metric does not belong to the character")
		}
	}

	// Check if the time tracking is already started
	timeTrackings, _ := r.TimeTrackingsRepo.GetTimeTrackingsByCharacterID(characterOID)
	for _, timeTracking := range timeTrackings {
		if timeTracking.EndTime.IsZero() {
			return nil, fmt.Errorf("focused session is already started")
		}
	}

	timeTracking := coredb.TimeTracking{
		ID:              primitive.NewObjectID(),
		CharacterID:     characterOID,
		CustomMetricID:  customMetricOID,
		StartTime:       time.Now(),
		MinDurationTime: 600,
		MaxDurationTime: 14400,
	}

	createdTimeTracking, err := r.TimeTrackingsRepo.CreateTimeTracking(&timeTracking)
	if err != nil {
		log.Printf("failed to insert time tracking: %v\n", err)
		return nil, err
	}

	return *createdTimeTracking, nil
}

func (r *TimeTrackingsResolver) UpdateTimeTracking(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	timeTrackingID := params.Args["id"].(string)
	timeTrackingOID, err := primitive.ObjectIDFromHex(timeTrackingID)
	if err != nil {
		return nil, err
	}

	timeTracking, err := r.TimeTrackingsRepo.GetTimeTrackingByID(timeTrackingOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get time tracking: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(timeTracking.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.UserID != user.ID {
		return nil, fmt.Errorf("permission denied")
	}

	if !timeTracking.EndTime.IsZero() {
		return nil, fmt.Errorf("focused session is already ended")
	}

	intDuration := params.Args["duration"].(int)
	duration := int32(intDuration)
	err = ValidateDuration(timeTracking.StartTime, duration)
	if err != nil {
		return nil, err
	}

	// JUST FOR TESTING
	// duration = 599 // Test for the min duration time
	// duration = 600 // Test for the min duration time
	// duration = 601 // Test for the min duration time
	// duration = 14400 // Test for the max duration time
	// duration = 14401 // Test for the max duration time

	if duration < timeTracking.MinDurationTime {
		duration = 0
		_, err := r.TimeTrackingsRepo.DeleteTimeTracking(timeTrackingOID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete time tracking: %v", err)
		}

		log.Printf("the period time is less than 10 min, so the time tracking will be deleted")
		return *timeTracking, nil
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

	return *timeTracking, nil
}
