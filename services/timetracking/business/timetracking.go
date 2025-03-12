package business

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/timetracking/entity"

	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"

	"github.com/samber/lo"
)

type TimeTrackingBusiness struct {
	coreClient     ICoreClient
	currencyClient ICurrencyClient
	notiClient     INotificationClient
	cache          ICache
}

func NewTimeTrackingsBusiness(coreClient ICoreClient, currencyClient ICurrencyClient, notiClient INotificationClient, cache ICache) *TimeTrackingBusiness {
	return &TimeTrackingBusiness{
		coreClient:     coreClient,
		currencyClient: currencyClient,
		notiClient:     notiClient,
		cache:          cache,
	}
}

func (biz *TimeTrackingBusiness) GetCurrentTimeTracking(ctx context.Context) (*entity.TimeTracking, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	currentTimeTrack, err := biz.cache.GetCurrentTimeTracking(ctx, authSession.ProfileID)
	if err == errors.ErrRedisNotFound {
		return nil, nil
	}

	return currentTimeTrack, err
}

func (biz *TimeTrackingBusiness) GetTotalCurrentTimeTracking(ctx context.Context, characterID string, timestamp time.Time) (int, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return 0, errors.ErrUnauthorized
	}

	// Check permissions
	authorized, err := biz.coreClient.CheckPermission(ctx, lo.ToPtr(authSession.ProfileID), lo.ToPtr(characterID), nil)
	if !authorized || err != nil {
		return 0, errors.ErrPermissionDenied
	}

	capturedRecord, err := biz.cache.GetCurrentCapturedRecord(ctx, authSession.ProfileID, characterID)
	if errors.HasCode(err, errors.ErrCodeRedisNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get current captured record: %v", err)
	}

	// Get the timetracking from the timestamp to now
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

func (biz *TimeTrackingBusiness) CreateTimeTracking(ctx context.Context, characterID string, categoryID *string, startTime time.Time) (*entity.TimeTracking, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	// Calculate the difference between the server time and the client time
	serverStartTime := time.Now()
	duration := serverStartTime.Sub(startTime)
	seconds := duration.Seconds()

	if seconds > utils.MaxTimeDifference {
		return nil, errors.NewGQLError(errors.ErrCodeOverMaxDifferenceDuration, "the period time is over the max difference duration")
	}

	authorized, err := biz.coreClient.CheckPermission(ctx, lo.ToPtr(authSession.ProfileID), lo.ToPtr(characterID), categoryID)
	if !authorized || err != nil {
		return nil, errors.ErrPermissionDenied
	}

	// Check if there is an active time tracking
	currentTimeTracking, err := biz.cache.GetCurrentTimeTracking(ctx, authSession.ProfileID)
	if err != nil && !errors.HasCode(err, errors.ErrCodeRedisNotFound) {
		return nil, err
	}

	if currentTimeTracking != nil {
		return nil, errors.NewGQLError(errors.ErrCodeTimeTrackingAlreadyExists, "time tracking already exists")
	}

	// Create a new time tracking
	timeTracking := &entity.TimeTracking{
		ID:          mongodb.GenObjectID(),
		CharacterID: characterID,
		StartTime:   startTime,
	}

	if categoryID != nil {
		timeTracking.CategoryID = categoryID
	}

	err = biz.cache.CreateTimeTracking(ctx, authSession.ProfileID, timeTracking)
	if err != nil {
		return nil, fmt.Errorf("failed to create time tracking: %v", err)
	}

	req := &entity.SendNotiReq{
		ProfileID: authSession.ProfileID,
		DeviceID:  authSession.DeviceID,
		Title:     "New Notification",
		Body:      "Start tracking!",
	}
	_, err = biz.notiClient.SendNotification(ctx, req)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
	} else {
		fmt.Println("Message sent successfully")
	}

	return timeTracking, nil
}

func (biz *TimeTrackingBusiness) UpdateTimeTracking(ctx context.Context) (*entity.TimeTracking, *entity.Fish, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, nil, errors.ErrUnauthorized
	}

	// Get the time tracking from Redis
	timeTracking, err := biz.cache.GetCurrentTimeTracking(ctx, authSession.ProfileID)
	if err == errors.ErrRedisNotFound {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	err = biz.cache.DeleteCurrentTimeTracking(ctx, authSession.ProfileID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete current time tracking: %v", err)
	}

	// Calculate the duration time
	endTime := time.Now()
	duration := int32(endTime.Sub(timeTracking.StartTime).Seconds())

	// Check if the duration time is in the valid range
	if duration < utils.MinDurationTime {
		return nil, nil, errors.NewGQLError(errors.ErrCodeUnderMinDuration, "the period time is less than min duration time")
	}

	if duration > utils.MaxDurationTime {
		duration = int32(utils.MaxDurationTime)
	}

	timeTracking.EndTime = timeTracking.StartTime.Add(time.Duration(duration) * time.Second)

	// Upsert the time tracking in captured record
	err = biz.cache.UpsertTimeTrackingInCapturedRecord(ctx, authSession.ProfileID, timeTracking, duration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to upsert time tracking in captured record: %v", err)
	}

	// Calculate the number of catches
	fishCatchingInterval := getFishCatchingInterval()
	numCatches := int(duration) / fishCatchingInterval

	updatedFish := &entity.Fish{
		ProfileID: authSession.ProfileID,
		Gold:      0,
		Normal:    0,
	}

	for i := 0; i < numCatches; i++ {
		catchResult, err := biz.currencyClient.CatchFish(ctx)
		if err != nil {
			return nil, nil, err
		}

		switch catchResult.FishType {
		case "Normal":
			updatedFish.Normal += catchResult.Number
		case "Gold":
			updatedFish.Gold += catchResult.Number
		}
	}

	err = biz.currencyClient.UpdateFish(ctx, updatedFish)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update fish %v", err)
	}

	req := &entity.SendNotiReq{
		ProfileID: authSession.ProfileID,
		DeviceID:  authSession.DeviceID,
		Title:     "New Notification",
		Body:      "Finish trackingg!",
	}
	_, err = biz.notiClient.SendNotification(ctx, req)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
	} else {
		fmt.Println("Message sent successfully")
	}

	return timeTracking, updatedFish, nil
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
