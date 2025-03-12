package business

import (
	"context"
	"time"

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
	GetCurrentCapturedRecord(ctx context.Context, profileID, characterID string) (*entity.CapturedRecord, error)
	UpsertTimeTrackingInCapturedRecord(ctx context.Context, profileID string, timeTracking *entity.TimeTracking, duration int32) error
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
