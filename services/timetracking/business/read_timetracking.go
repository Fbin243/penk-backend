package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/services/timetracking/entity"

	rdb "tenkhours/pkg/db/redis"

	"github.com/samber/lo"
)

type TimeTrackingBusiness struct {
	coreClient       ICoreClient
	currencyClient   ICurrencyClient
	notiClient       INotificationClient
	cache            ICache
	timetrackingRepo ITimeTrackingRepo
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

	timetrackings, err := biz.timetrackingRepo.FindByCharacterID(ctx, characterID)
	if err != nil {
		return 0, err
	}

	totalTime := 0
	for _, timetracking := range timetrackings {
		if timetracking.StartTime.Before(timestamp) && timetracking.EndTime.After(timestamp) {
			totalTime += int(timetracking.EndTime.Sub(timestamp).Seconds())
		}
		if timestamp.Before(timetracking.StartTime) {
			totalTime += int(timetracking.EndTime.Sub(timetracking.StartTime).Seconds())
		}
	}

	return totalTime, nil
}
