package business

// func (biz *TimeTrackingBusiness) GetCurrentTimeTracking(ctx context.Context) (*entity.TimeTracking, error) {
// 	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
// 	if !ok {
// 		return nil, errors.ErrUnauthorized
// 	}

// 	currentTimeTrack, err := biz.cache.GetCurrentTimeTracking(ctx, authSession.ProfileID)
// 	if err == errors.ErrRedisNotFound {
// 		return nil, nil
// 	}

// 	return currentTimeTrack, err
// }

// func (biz *TimeTrackingBusiness) GetTotalCurrentTimeTracking(ctx context.Context, timestamp time.Time) (int, error) {
// 	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
// 	if !ok {
// 		return 0, errors.ErrUnauthorized
// 	}

// 	timetrackings, err := biz.timetrackingRepo.FindByCharacterID(ctx, authSession.CurrentCharacterID)
// 	if err != nil {
// 		return 0, err
// 	}

// 	totalTime := 0
// 	for _, timetracking := range timetrackings {
// 		if timetracking.StartTime.Before(timestamp) && timetracking.EndTime.After(timestamp) {
// 			totalTime += int(timetracking.EndTime.Sub(timestamp).Seconds())
// 		}
// 		if timestamp.Before(timetracking.StartTime) {
// 			totalTime += int(timetracking.EndTime.Sub(timetracking.StartTime).Seconds())
// 		}
// 	}

// 	return totalTime, nil
// }
