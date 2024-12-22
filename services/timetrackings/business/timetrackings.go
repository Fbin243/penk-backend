package business

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/analytics/graph/model"
	"tenkhours/services/core/repo"
	fishBiz "tenkhours/services/currency/business"
	fishRepo "tenkhours/services/currency/repo"
	timetrackingsRepo "tenkhours/services/timetrackings/repo"
	"tenkhours/services/timetrackings/repo/rpc"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeTrackingsBusiness struct {
	TimeTrackingsRepo *timetrackingsRepo.TimeTrackingsRepo
	RPCClient         *rpc.RPCClient
	FishRepo          *fishRepo.FishRepo
	ProfilesRepo      *repo.ProfilesRepo
	RedisClient       *redis.Client
}

func NewTimeTrackingsBusiness(timeTrackingsRepo *timetrackingsRepo.TimeTrackingsRepo, coreClient *rpc.RPCClient, fishRepo *fishRepo.FishRepo, profilesRepo *repo.ProfilesRepo, redisClient *redis.Client) *TimeTrackingsBusiness {
	return &TimeTrackingsBusiness{
		TimeTrackingsRepo: timeTrackingsRepo,
		RPCClient:         coreClient,
		FishRepo:          fishRepo,
		ProfilesRepo:      profilesRepo,
		RedisClient:       redisClient,
	}
}

func (biz *TimeTrackingsBusiness) GetCurrentTimeTracking(ctx context.Context) (*timetrackingsRepo.TimeTracking, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	// Get the current time tracking from Redis
	currentTimetrackingJSON, err := biz.RedisClient.Get(ctx, db.GetTimeTrackingKey(authSession.ProfileID.Hex())).Result()
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

func (biz *TimeTrackingsBusiness) GetTotalCurrentTimeTracking(ctx context.Context, characterID primitive.ObjectID, timestamp time.Time) (int, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return 0, errors.Unauthorized()
	}

	// Check permissions
	authorized, err := biz.RPCClient.CheckPermission(ctx, authSession.ProfileID.Hex(), characterID.Hex(), primitive.NilObjectID.Hex())
	if !authorized || err != nil {
		return 0, errors.PermissionDenied()
	}

	// Get the timetrackings from the current captured record in Redis
	capturedRecordJSON, err := biz.RedisClient.HGet(ctx, db.GetCapturedRecordKey(authSession.ProfileID.Hex()), characterID.Hex()).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("failed to get captured record from redis: %v", err)
	}

	var capturedRecord model.CapturedRecord
	err = json.Unmarshal([]byte(capturedRecordJSON), &capturedRecord)
	if err != nil {
		return 0, fmt.Errorf("failed to deserialize captured record: %v", err)
	}

	// Get the timetrackings from the timestamp to now
	totalTime := 0
	for _, timeTracking := range capturedRecord.TimeTrackings {
		if timestamp.After(timeTracking.StartTime) && timestamp.Before(timeTracking.EndTime) {
			// Case 1: The timestamp after the startTime and before the endTime
			totalTime += int(timeTracking.EndTime.Sub(timestamp).Seconds())
		} else if timestamp.Before(timeTracking.StartTime) {
			// Case 2: The timestamp before the startTime
			totalTime += int(timeTracking.Time)
		}
	}

	return totalTime, nil
}

// Helper function to get time interval fish catching
func getFishCatchingInterval() int {
	fishCatchingIntervalStr := os.Getenv("FISH_CATCHING_INTERVAL_SECONDS")
	fishCatchingInterval := 5 // Default value (5 seconds) for testing

	if fishCatchingIntervalStr != "" {
		interval, err := strconv.Atoi(fishCatchingIntervalStr) // convert string to int
		if err != nil {
			log.Printf("Invalid FISH_CATCHING_INTERVAL_SECONDS: %v, using default 5 seconds", err)
		} else {
			fishCatchingInterval = interval
		}
	}
	return fishCatchingInterval
}

func (biz *TimeTrackingsBusiness) CreateTimeTracking(ctx context.Context, characterID primitive.ObjectID, metricID *primitive.ObjectID, startTime time.Time) (*timetrackingsRepo.TimeTracking, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	// Calculate the difference between the server time and the client time
	serverStartTime := time.Now()
	duration := serverStartTime.Sub(startTime)
	seconds := duration.Seconds()

	if seconds > utils.MaxTimeDifference {
		return nil, errors.NewError(errors.ErrCodeOverMaxDifferenceDuration, "the period time is over the max difference duration")
	}

	var metricOID primitive.ObjectID
	if metricID != nil {
		metricOID = *metricID
	}

	authorized, err := biz.RPCClient.CheckPermission(ctx, authSession.ProfileID.Hex(), characterID.Hex(), metricOID.Hex())
	if !authorized || err != nil {
		return nil, errors.PermissionDenied()
	}

	// Check if there is an active time tracking
	currentTimeTracking, err := biz.GetCurrentTimeTracking(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current time tracking: %v", err)
	}

	if currentTimeTracking != nil {
		return nil, errors.NewError(errors.ErrCodeTimeTrackingAlreadyExists, "time tracking already exists")
	}

	// Create a new time tracking
	timeTracking := timetrackingsRepo.TimeTracking{
		ID:              primitive.NewObjectID(),
		CharacterID:     characterID,
		StartTime:       startTime,
		MinDurationTime: utils.MinDurationTime,
		MaxDurationTime: utils.MaxDurationTime,
		CustomMetricID:  metricOID,
	}

	// Save the time tracking to Redis
	timeTrackingJSON, err := json.Marshal(timeTracking)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize time tracking: %v", err)
	}

	err = biz.RedisClient.Set(ctx, db.GetTimeTrackingKey(authSession.ProfileID.Hex()), timeTrackingJSON, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save time tracking to redis: %v", err)
	}

	return &timeTracking, nil
}

func (biz *TimeTrackingsBusiness) UpdateTimeTracking(ctx context.Context) (*timetrackingsRepo.TimeTracking, *fishRepo.Fish, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, nil, errors.Unauthorized()
	}

	// Get the time tracking from Redis
	profileID := authSession.ProfileID.Hex()
	val, err := biz.RedisClient.Get(ctx, db.GetTimeTrackingKey(profileID)).Result()
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
	if duration < timeTracking.MinDurationTime {
		duration = 0
		return nil, nil, errors.NewError(errors.ErrCodeUnderMinDuration, "the period time is less than min duration time")
	}

	if duration > timeTracking.MaxDurationTime {
		duration = int32(timeTracking.MaxDurationTime)
	}

	timeTracking.EndTime = timeTracking.StartTime.Add(time.Duration(duration) * time.Second)

	// Check if the captured record already exists in Redis
	capturedRecord := model.CapturedRecord{}
	capturedRecordJSON, err := biz.RedisClient.HGet(ctx, db.GetCapturedRecordKey(profileID), timeTracking.CharacterID.Hex()).Result()
	if err == redis.Nil {
		// Make a new captured record if it doesn't exist
		capturedRecord = model.CapturedRecord{
			ID:               primitive.NewObjectID(),
			Timestamp:        utils.ResetTimeToBeginningOfDay(timeTracking.EndTime),
			TotalFocusedTime: 0,
			Metadata: model.CapturedRecordMetadata{
				CharacterID: timeTracking.CharacterID,
				ProfileID:   authSession.ProfileID,
			},
			TimeTrackings: []model.CapturedRecordTimeTracking{},
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

	// Add the time tracking to the captured record
	capturedRecord.TimeTrackings = append(capturedRecord.TimeTrackings, model.CapturedRecordTimeTracking{
		CustomMetricID: timeTracking.CustomMetricID,
		Time:           int32(duration),
		StartTime:      timeTracking.StartTime,
		EndTime:        timeTracking.EndTime,
	})

	if !timeTracking.CustomMetricID.IsZero() {
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
	}

	// Update the time in the character by gRPC
	success, err := biz.RPCClient.UpdateTimeInCharacter(ctx, timeTracking.CharacterID.Hex(), timeTracking.CustomMetricID.Hex(), duration)
	if !success || err != nil {
		return nil, nil, err
	}

	// Upsert the captured record to Redis
	capturedRecordBytes, err := json.Marshal(capturedRecord)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize captured record: %v", err)
	}

	err = biz.RedisClient.HSet(ctx, db.GetCapturedRecordKey(profileID), timeTracking.CharacterID.Hex(), capturedRecordBytes).Err()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save captured record to redis: %v", err)
	}

	// Calculate the number of catches
	fishCatchingInterval := getFishCatchingInterval()
	numCatches := int(duration) / fishCatchingInterval

	updatedFish := &fishRepo.Fish{
		ProfileID: authSession.ProfileID,
		Gold:      0,
		Normal:    0,
	}

	fishBiz := fishBiz.NewFishBusiness(biz.FishRepo, nil, biz.ProfilesRepo, biz.RedisClient)

	for i := 0; i < numCatches; i++ {
		catchResult, err := fishBiz.CatchFish(ctx, authSession.ProfileID)
		if err != nil {
			return nil, nil, err
		}

		switch catchResult.FishType {
		case "normal":
			updatedFish.Normal += catchResult.Number
		case "gold":
			updatedFish.Gold += catchResult.Number
		}
	}

	_, err = fishBiz.UpdateFishFromFinishSession(updatedFish, authSession.ProfileID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update fish %v", err)
	}

	return &timeTracking, updatedFish, nil
}
