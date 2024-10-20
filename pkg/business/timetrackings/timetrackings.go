package timetrackings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/db/timetrackingsdb"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimeTrackingsHandler struct {
	TimeTrackingsRepo *timetrackingsdb.TimeTrackingsRepo
	CharactersRepo    *coredb.CharactersRepo
	RedisClient       *redis.Client
}

func NewTimeTrackingsHandler(timeTrackingsRepo *timetrackingsdb.TimeTrackingsRepo, charactersRepo *coredb.CharactersRepo, redisClient *redis.Client) *TimeTrackingsHandler {
	return &TimeTrackingsHandler{
		TimeTrackingsRepo: timeTrackingsRepo,
		CharactersRepo:    charactersRepo,
		RedisClient:       redisClient,
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
		StartTime:       startTime,
		MinDurationTime: 600,
		MaxDurationTime: 14400,
	}

	if metricID != nil {
		timeTracking.CustomMetricID = *metricID
	}

	timeTrackingJSON, err := json.Marshal(timeTracking)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize time tracking: %v", err)
	}

	err = r.RedisClient.Set(ctx, timeTracking.ID.Hex(), timeTrackingJSON, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save time tracking to redis: %v", err)
	}

	return &timeTracking, nil
}

func (r *TimeTrackingsHandler) UpdateTimeTracking(ctx context.Context, id primitive.ObjectID, characterID primitive.ObjectID, metricID *primitive.ObjectID) (*timetrackingsdb.TimeTracking, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	endTime := time.Now()

	val, err := r.RedisClient.Get(ctx, id.Hex()).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("time tracking not found in redis")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get time tracking from redis: %v", err)
	}

	var timeTracking timetrackingsdb.TimeTracking
	err = json.Unmarshal([]byte(val), &timeTracking)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize time tracking: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
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

	if duration < timeTracking.MinDurationTime {
		duration = 0
		err = r.RedisClient.Del(ctx, id.Hex()).Err()
		if err != nil {
			return nil, fmt.Errorf("failed to delete time tracking from redis: %v", err)
		}

		log.Printf("the period time is less than 10 min, so the time tracking will be deleted")
		return &timeTracking, nil
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

	err = r.RedisClient.Set(ctx, timeTracking.ID.Hex(), timeTracking, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to update time tracking in redis: %v", err)
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return &timeTracking, nil
}
