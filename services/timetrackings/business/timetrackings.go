package business

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/services/analytics/graph/model"
	"tenkhours/services/core/repo"
	fishBiz "tenkhours/services/currency/business"
	fishRepo "tenkhours/services/currency/repo"
	timetrackingsRepo "tenkhours/services/timetrackings/repo"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeTrackingsBusiness struct {
	TimeTrackingsRepo *timetrackingsRepo.TimeTrackingsRepo
	CharactersRepo    *repo.CharactersRepo
	FishRepo          *fishRepo.FishRepo
	RedisClient       *redis.Client
}

func NewTimeTrackingsBusiness(timeTrackingsRepo *timetrackingsRepo.TimeTrackingsRepo, charactersRepo *repo.CharactersRepo, fishRepo *fishRepo.FishRepo, redisClient *redis.Client) *TimeTrackingsBusiness {
	return &TimeTrackingsBusiness{
		TimeTrackingsRepo: timeTrackingsRepo,
		CharactersRepo:    charactersRepo,
		FishRepo:          fishRepo,
		RedisClient:       redisClient,
	}
}

func (biz *TimeTrackingsBusiness) GetCurrentTimeTracking(ctx context.Context) (*timetrackingsRepo.TimeTracking, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// Get the current time tracking from Redis
	currentTimetrackingJSON, err := biz.RedisClient.Get(ctx, profile.ID.Hex()).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get time tracking from redis: %v", err)
	}

	var currentTimetracking timetrackingsRepo.TimeTracking
	err = json.Unmarshal([]byte(currentTimetrackingJSON), &currentTimetracking)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize time tracking: %v", err)
	}

	return &currentTimetracking, nil
}

func (biz *TimeTrackingsBusiness) CreateTimeTracking(ctx context.Context, characterID primitive.ObjectID, metricID *primitive.ObjectID, startTime time.Time) (*timetrackingsRepo.TimeTracking, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// Calculate the difference between the server time and the client time
	// serverStartTime := time.Now()
	// duration := serverStartTime.Sub(startTime)
	// seconds := duration.Seconds()

	// if seconds > 20 {
	// 	return nil, fmt.Errorf("server timeout, failed to start a new session")
	// }

	//for testing only
	startTime = time.Now()
	// Check permissions
	character, err := biz.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
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
			return nil, errors.ErrorPermissionDenied
		}
	}

	// Check if there is an active time tracking
	currentTimeTracking, err := biz.GetCurrentTimeTracking(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current time tracking: %v", err)
	}

	if currentTimeTracking != nil {
		return nil, fmt.Errorf("there is an active time tracking")
	}

	// Create a new time tracking
	timeTracking := timetrackingsRepo.TimeTracking{
		ID:              primitive.NewObjectID(),
		CharacterID:     characterID,
		StartTime:       startTime,
		MinDurationTime: 600,
		MaxDurationTime: 14400,
	}

	if metricID != nil {
		timeTracking.CustomMetricID = *metricID
	}

	// Save the time tracking to Redis
	timeTrackingJSON, err := json.Marshal(timeTracking)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize time tracking: %v", err)
	}

	err = biz.RedisClient.Set(ctx, profile.ID.Hex(), timeTrackingJSON, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save time tracking to redis: %v", err)
	}

	// Create null data for the first time
	fishKey := fmt.Sprintf("fish:%s", profile.ID.Hex())
	fish := fishRepo.Fish{
		ID:        primitive.NewObjectID(),
		ProfileID: profile.ID,
		Normal:    0,
		Gold:      0,
	}

	fishJSON, err := json.Marshal(fish)
	err = biz.RedisClient.Set(ctx, fishKey, fishJSON, 4*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save initial fish data to Redis: %v", err)
	}

	// Start fish-catching goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		log.Println("Start the goroutine")

		redisCtx := context.Background() // tạo ctx con do ctx cha bị huỷ sớm

		for {
			select {
			case <-ticker.C:
				fishBiz := &fishBiz.FishBusiness{RedisClient: biz.RedisClient}
				fishType, err := fishBiz.CatchFish(redisCtx, profile.ID)
				if err != nil {
					log.Printf("Failed to catch fish: %v", err)
					return
				}

				log.Printf("Caught fish: %s", fishType)
			}

			// Check if `fishData` still exists in Redis
			if _, err := biz.RedisClient.Get(ctx, fishKey).Result(); err == redis.Nil {
				fmt.Println("Fish data not found, stopping goroutine.")
				return
			}
		}
	}()

	return &timeTracking, nil
}

func (biz *TimeTrackingsBusiness) UpdateTimeTracking(ctx context.Context) (*timetrackingsRepo.TimeTracking, *fishRepo.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, nil, errors.ErrorUnauthorized
	}

	// Get the time tracking from Redis
	profileID := profile.ID.Hex()
	val, err := biz.RedisClient.Get(ctx, profileID).Result()
	if err == redis.Nil {
		return nil, nil, fmt.Errorf("time tracking not found in redis")
	} else if err != nil {
		return nil, nil, fmt.Errorf("failed to get time tracking from redis: %v", err)
	}

	var timeTracking timetrackingsRepo.TimeTracking
	err = json.Unmarshal([]byte(val), &timeTracking)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize time tracking: %v", err)
	}

	// Delete the current time tracking from Redis
	err = biz.RedisClient.Del(ctx, profileID).Err()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete time tracking from redis: %v", err)
	}

	// Calculate the duration time
	endTime := time.Now()
	duration := int32(endTime.Sub(timeTracking.StartTime).Seconds())
	// TODO: TESTING
	// duration = 599 // Test for the min duration time
	// duration = 600 // Test for the min duration time
	// duration = 601 // Test for the min duration time
	// duration = 14400 // Test for the max duration time
	// duration = 14401 // Test for the max duration time

	// Check if the duration time is in the valid range
	// if duration < timeTracking.MinDurationTime {
	// 	duration = 0
	// 	log.Printf("the period time is less than 10 min, so the time tracking will be deleted")
	// 	return &timeTracking, nil
	// }

	if duration > timeTracking.MaxDurationTime {
		duration = int32(timeTracking.MaxDurationTime)
		log.Printf("the period time is more than 4 hours, so the time tracking will be limited to 4 hours")
	}

	timeTracking.EndTime = timeTracking.StartTime.Add(time.Duration(duration) * time.Second)

	// Check if the captured record already exists in Redis
	capturedRecord := model.CapturedRecord{}
	capturedRecordJSON, err := biz.RedisClient.HGet(ctx, db.CapturedRecordKey+profile.ID.Hex(), timeTracking.CharacterID.Hex()).Result()
	if err == redis.Nil {
		// Make a new captured record if it doesn't exist
		capturedRecord = model.CapturedRecord{
			ID:               primitive.NewObjectID(),
			Timestamp:        timeTracking.StartTime,
			TotalFocusedTime: 0,
			Metadata: model.CapturedRecordMetadata{
				CharacterID: timeTracking.CharacterID,
				ProfileID:   profile.ID,
			},
			CustomMetrics: []model.CapturedRecordCustomMetric{},
		}
	} else if err != nil {
		return nil, nil, fmt.Errorf("failed to get captured record from redis: %v", err)
	} else {
		err = json.Unmarshal([]byte(capturedRecordJSON), &capturedRecord)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to deserialize captured record: %v", err)
		}
	}

	// Add the time to the captured record for character
	capturedRecord.TotalFocusedTime += duration

	// Get the character to update the time
	character, err := biz.CharactersRepo.GetCharacterByID(timeTracking.CharacterID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get character: %v", err)
	}

	character.TotalFocusedTime += duration
	if !timeTracking.CustomMetricID.IsZero() {
		for i, customMetric := range character.CustomMetrics {
			if customMetric.ID == timeTracking.CustomMetricID {
				character.CustomMetrics[i].Time += int32(duration)

				// Check if this custom metric already exists in the captured record
				found := false
				for j, capturedCustomMetric := range capturedRecord.CustomMetrics {
					// If it exists, add the time to it
					if capturedCustomMetric.ID == timeTracking.CustomMetricID {
						capturedRecord.CustomMetrics[j].Time += int32(duration)
						found = true
					}
				}

				// If it doesn't exist, create a new one
				if !found {
					capturedRecord.CustomMetrics = append(capturedRecord.CustomMetrics, model.CapturedRecordCustomMetric{
						ID:   timeTracking.CustomMetricID,
						Time: int32(duration),
					})
				}
				break
			}
		}
	}

	_, err = biz.TimeTrackingsRepo.CreateTimeTracking(&timeTracking)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create time tracking in DB: %v", err)
	}

	_, err = biz.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update character: %v", err)
	}

	// Upsert the captured record to Redis
	capturedRecordBytes, err := json.Marshal(capturedRecord)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize captured record: %v", err)
	}

	err = biz.RedisClient.HSet(ctx, db.CapturedRecordKey+profile.ID.Hex(), character.ID.Hex(), capturedRecordBytes).Err()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save captured record to redis: %v", err)
	}

	// retrieve fish data from Redis and delete cache
	updatedFish := &fishRepo.Fish{
		ProfileID: profile.ID,
		Gold:      0,
		Normal:    0,
	}
	fishData, err := biz.RedisClient.Get(ctx, fmt.Sprintf("fish:%s", profile.ID.Hex())).Result()
	if err == redis.Nil {
		log.Printf("No fish data found for profile %s", profileID)
	} else if err != nil {
		return nil, nil, fmt.Errorf("failed to get fish data from redis: %v", err)
	} else {
		// Delete the fish data cache
		err = biz.RedisClient.Del(ctx, fmt.Sprintf("fish:%s", profile.ID.Hex())).Err()
		if err != nil {
			log.Printf("failed to delete fish data from redis: %v", err)
		}

		err = json.Unmarshal([]byte(fishData), updatedFish)
		if err != nil {
			log.Printf("failed to unmarshal fish data from redis: %v", err)
		}

		fishBiz := &fishBiz.FishBusiness{FishRepo: biz.FishRepo}

		_, err = fishBiz.UpdateFishFromRedis(updatedFish, profile.ID)

		log.Printf("Successfully deleted fish data cache for profile %s", profileID)
	}

	return &timeTracking, updatedFish, nil
}
