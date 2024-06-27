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

func NewTimeTrackingsResolver() *TimeTrackingsResolver {
	return &TimeTrackingsResolver{
		TimeTrackingsRepo: coredb.NewTimeTrackingsRepo(),
		CharactersRepo:    coredb.NewCharactersRepo(),
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
		ID:             primitive.NewObjectID(),
		CharacterID:    characterOID,
		CustomMetricID: customMetricOID,
		StartTime:      time.Now(),
	}

	insertResult, err := r.TimeTrackingsRepo.CreateTimeTracking(timeTracking)
	if err != nil {
		log.Printf("failed to insert time tracking: %v\n", err)
		return nil, err
	}

	return insertResult, nil
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

	character.TotalFocusedTime += int32(duration)
	if !timeTracking.CustomMetricID.IsZero() {
		for _, customMetric := range character.CustomMetrics {
			if customMetric.ID == timeTracking.CustomMetricID {
				customMetric.Time += int32(duration)
				break
			}
		}
	}

	updateResult, err := r.TimeTrackingsRepo.UpdateTimeTracking(timeTracking)
	if err != nil {
		return nil, fmt.Errorf("failed to update time tracking: %v", err)
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return updateResult, nil
}
