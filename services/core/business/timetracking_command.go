package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (biz *TimeTrackingBusiness) UpsertTimeTracking(ctx context.Context, input *entity.TimeTrackingInput) (*entity.TimeTracking, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the user has permission to access the reference ID
	err = biz.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   input.ReferenceID,
			Type: input.ReferenceType,
		},
	})
	if err != nil {
		return nil, err
	}

	// Check if there is a time tracking in timestamp
	currentTimeTracking, err := biz.timetrackingRepo.FindByReferenceIDAndTimestamp(ctx, input.ReferenceID, input.Timestamp)

	if err == errors.ErrMongoNotFound {
		// No time tracking found, create a new one

		// Get category id
		var categoryID *string
		switch input.ReferenceType {
		case entity.EntityTypeTask:
			task, err := biz.taskRepo.FindByID(ctx, input.ReferenceID)
			if err != nil {
				return nil, err
			}
			categoryID = task.CategoryID
		case entity.EntityTypeHabit:
			habit, err := biz.habitRepo.FindByID(ctx, input.ReferenceID)
			if err != nil {
				return nil, err
			}
			categoryID = habit.CategoryID
		}

		timeTracking := &entity.TimeTracking{
			BaseEntity:    &base.BaseEntity{},
			CharacterID:   authSession.CurrentCharacterID,
			CategoryID:    categoryID,
			ReferenceID:   lo.ToPtr(input.ReferenceID),
			ReferenceType: lo.ToPtr(input.ReferenceType),
			Timestamp:     input.Timestamp,
			Duration:      input.Duration,
		}

		return biz.timetrackingRepo.InsertOne(ctx, timeTracking)
	} else if err != nil {
		return nil, err
	}

	// Time tracking found, update the existing one
	currentTimeTracking.Duration += input.Duration
	return biz.timetrackingRepo.FindAndUpdateByID(ctx, currentTimeTracking.ID, currentTimeTracking)
}

// func (biz *TimeTrackingBusiness) CreateTimeTracking(ctx context.Context, input *entity.TimeTrackingInput) (*entity.TimeTracking, error) {
// 	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
// 	if !ok {
// 		return nil, errors.ErrUnauthorized
// 	}

// 	// Calculate the difference between the server time and the client time
// 	serverStartTime := time.Now()
// 	if serverStartTime.Before(input.StartTime) {
// 		return nil, errors.NewGQLError(errors.ErrCodeBadRequest, "start time is in the future")
// 	}

// 	duration := serverStartTime.Sub(input.StartTime)
// 	seconds := duration.Seconds()

// 	if seconds > utils.MaxTimeDifference {
// 		return nil, errors.NewGQLError(errors.ErrCodeOverMaxDifferenceDuration, "the period time is over the max difference duration")
// 	}

// 	permEntities := []PermissionEntity{}
// 	if input.CategoryID != nil {
// 		permEntities = append(permEntities, PermissionEntity{
// 			ID:   *input.CategoryID,
// 			Type: entity.EntityTypeCategory,
// 		})
// 	}

// 	if input.ReferenceID != nil {
// 		permEntities = append(permEntities, PermissionEntity{
// 			ID:   *input.ReferenceID,
// 			Type: *input.ReferenceType,
// 		})
// 	}

// 	err := biz.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, permEntities)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check if there is an active time tracking
// 	currentTimeTracking, err := biz.cache.GetCurrentTimeTracking(ctx, authSession.ProfileID)
// 	if err != nil && !errors.HasCode(err, errors.ErrCodeRedisNotFound) {
// 		return nil, err
// 	}

// 	if currentTimeTracking != nil {
// 		return nil, errors.NewGQLError(errors.ErrCodeTimeTrackingAlreadyExists, "time tracking already exists")
// 	}

// 	// Create a new time tracking
// 	timeTracking := &entity.TimeTracking{
// 		BaseEntity: &base.BaseEntity{
// 			ID: mongodb.GenObjectID(),
// 		},
// 		CharacterID:   authSession.CurrentCharacterID,
// 		CategoryID:    input.CategoryID,
// 		ReferenceID:   input.ReferenceID,
// 		ReferenceType: input.ReferenceType,
// 		StartTime:     input.StartTime,
// 	}

// 	err = biz.cache.CreateTimeTracking(ctx, authSession.ProfileID, timeTracking)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create time tracking: %v", err)
// 	}

// 	// Send noti
// 	req := &entity.SendNotiReq{
// 		ProfileID: authSession.ProfileID,
// 		DeviceID:  authSession.DeviceID,
// 		Title:     "New Notification",
// 		Body:      "Start tracking!",
// 	}
// 	_, err = biz.notiClient.SendNotification(ctx, req)
// 	if err != nil {
// 		log.Printf("Failed to send notification: %v", err)
// 	} else {
// 		fmt.Println("Message sent successfully")
// 	}

// 	return timeTracking, nil
// }

// func (biz *TimeTrackingBusiness) UpdateTimeTracking(ctx context.Context) (*entity.TimeTracking, error) {
// 	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
// 	if !ok {
// 		return nil, errors.ErrUnauthorized
// 	}

// 	// Get the time tracking from Redis
// 	timeTracking, err := biz.cache.GetCurrentTimeTracking(ctx, authSession.ProfileID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = biz.cache.DeleteCurrentTimeTracking(ctx, authSession.ProfileID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to delete current time tracking: %v", err)
// 	}

// 	// Calculate the duration time
// 	endTime := time.Now()
// 	duration := int32(endTime.Sub(timeTracking.StartTime).Seconds())

// 	// Check if the duration time is in the valid range
// 	if duration < utils.MinDurationTime {
// 		return nil, errors.NewGQLError(errors.ErrCodeUnderMinDuration, "the period time is less than min duration time")
// 	}

// 	if duration > utils.MaxDurationTime {
// 		duration = int32(utils.MaxDurationTime)
// 	}

// 	timeTracking.EndTime = timeTracking.StartTime.Add(time.Duration(duration) * time.Second)
// 	timeTracking, err = biz.timetrackingRepo.InsertOne(ctx, timeTracking)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create habit log for completion type time
// 	if timeTracking.ReferenceType != nil &&
// 		*timeTracking.ReferenceType == entity.EntityTypeHabit {
// 		habit, err := biz.habitRepo.FindByID(ctx, *timeTracking.ReferenceID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if habit.CompletionType == entity.CompletionTypeTime {
// 			err := biz.habitLogRepo.UpsertByTimestamp(ctx, timeTracking.EndTime, &entity.HabitLog{
// 				BaseEntity: &base.BaseEntity{},
// 				HabitID:    habit.ID,
// 				Timestamp:  timeTracking.EndTime,
// 				Value:      float64(duration),
// 			})
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}

// 	// Send noti
// 	req := &entity.SendNotiReq{
// 		ProfileID: authSession.ProfileID,
// 		DeviceID:  authSession.DeviceID,
// 		Title:     "New Notification",
// 		Body:      "Finish trackingg!",
// 	}
// 	_, err = biz.notiClient.SendNotification(ctx, req)
// 	if err != nil {
// 		log.Printf("Failed to send notification: %v", err)
// 	} else {
// 		fmt.Println("Message sent successfully")
// 	}

// 	return timeTracking, nil
// }
