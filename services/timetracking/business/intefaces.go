package business

import (
	"context"
	"time"

	"tenkhours/pkg/db/base"
	"tenkhours/services/timetracking/entity"
)

type ITimeTrackingBusiness interface {
	GetCurrentTimeTracking(ctx context.Context) (*entity.TimeTracking, error)
	GetTotalCurrentTimeTracking(ctx context.Context, characterID string, timestamp time.Time) (int, error)
	CreateTimeTracking(ctx context.Context, characterID string, categoryID *string, startTime time.Time) (*entity.TimeTracking, error)
	UpdateTimeTracking(ctx context.Context) (*entity.TimeTracking, *entity.Fish, error)
}

type ICache interface {
	GetCurrentTimeTracking(ctx context.Context, profileID string) (*entity.TimeTracking, error)
	DeleteCurrentTimeTracking(ctx context.Context, profileID string) error
	CreateTimeTracking(ctx context.Context, profileID string, timeTracking *entity.TimeTracking) error
}

type ICoreClient interface {
	CheckPermission(ctx context.Context, profileID, characterID, categoryID *string) (bool, error)
}

type ICurrencyClient interface {
	CatchFish(ctx context.Context) (*entity.CatchFishResult, error)
	UpdateFish(ctx context.Context, fish *entity.Fish) error
}

type INotificationClient interface {
	SendNotification(ctx context.Context, req *entity.SendNotiReq) (bool, error)
}

type ITimeTrackingRepo interface {
	base.IBaseRepo[entity.TimeTracking]
	FindByTimestamp(ctx context.Context, timestamp time.Time) ([]entity.TimeTracking, error)
	FindByCategoryID(ctx context.Context, categoryID string) ([]entity.TimeTracking, error)
	FindByCharacterID(ctx context.Context, characterID string) ([]entity.TimeTracking, error)
}
