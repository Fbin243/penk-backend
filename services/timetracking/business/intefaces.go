package business

import (
	"context"

	"tenkhours/services/timetracking/entity"
)

type ITimeTrackingBusiness interface{}

type ICache interface {
	GetCurrentTimeTracking(ctx context.Context, profileID string) (*entity.TimeTracking, error)
	DeleteCurrentTimeTracking(ctx context.Context, profileID string) error
	CreateTimeTracking(ctx context.Context, profileID string, timeTracking *entity.TimeTracking) error
	GetCurrentCapturedRecord(ctx context.Context, profileID, characterID string) (*entity.CapturedRecord, error)
	UpsertTimeTrackingInCapturedRecord(ctx context.Context, profileID string, timeTracking *entity.TimeTracking, duration int32) error
}

type ICoreClient interface {
	UpdateTimeInCharacter(ctx context.Context, characterID string, metricID *string, time int32) (bool, error)
	CheckPermission(ctx context.Context, profileID, characterID string, metricID *string) (bool, error)
}

type ICurrencyClient interface {
	CatchFish(ctx context.Context) (*entity.CatchFishResult, error)
	UpdateFish(ctx context.Context, fish *entity.Fish) error
}
