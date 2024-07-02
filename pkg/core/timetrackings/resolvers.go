package timetrackings

import (
	"fmt"
	"log"
	"time"

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
	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, err
	}

	customMetricID, ok := params.Args["customMetricID"].(string)
	customMetricOID := primitive.ObjectID{}
	if ok {
		customMetricOID, err = primitive.ObjectIDFromHex(customMetricID)
		if err != nil {
			return nil, err
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

	_, err = r.TimeTrackingsRepo.CreateTimeTracking(timeTracking)
	if err != nil {
		log.Printf("failed to insert time tracking: %v\n", err)
		return nil, err
	}

	return timeTracking.ID.Hex(), nil
}

func (r *TimeTrackingsResolver) UpdateTimeTracking(params graphql.ResolveParams) (interface{}, error) {
	timeTrackingID := params.Args["id"].(string)
	timeTrackingOID, err := primitive.ObjectIDFromHex(timeTrackingID)
	if err != nil {
		return nil, err
	}

	timeTracking, err := r.TimeTrackingsRepo.GetTimeTrackingByID(timeTrackingOID)
	if err != nil {
		return nil, fmt.Errorf("time tracking not found: %v", err)
	}

	timeTracking.EndTime = time.Now()

	character, err := r.CharactersRepo.GetCharacterByID(timeTracking.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	duration := timeTracking.EndTime.Sub(timeTracking.StartTime).Seconds()
	if int32(duration) < timeTracking.MinDurationTime {
		duration = 0
		return nil, fmt.Errorf("the period time is less than 10 min")
	}

	if int32(duration) > timeTracking.MaxDurationTime {
		duration = float64(timeTracking.MaxDurationTime)
		return nil, fmt.Errorf("the period time is more than 4 hours, so the period time will set to 4 hours")
	}

	character.TotalFocusedTime += int32(duration)
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

	return timeTrackingID, nil
}
