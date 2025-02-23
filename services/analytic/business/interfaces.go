package business

import (
	"context"
	"time"

	"tenkhours/services/analytic/entity"
)

// Business
type IAnalyticBusiness interface {
	GetAnalyticResults(ctx context.Context, characterID *string, startTime, endTime *time.Time, analyticSections []entity.AnalyticSection) (map[string]interface{}, error)
	DeleteCapturedRecords(ctx context.Context, profileID string) error
}

// Repository
type ICapturedRecordRepo interface {
	GetCapturedRecords(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error)
	DeleteCapturedRecords(ctx context.Context, profileID string) error
}

type ICache interface {
	GetCapturedRecords(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error)
}

// RPC client
type ICoreClient interface {
	CheckPermission(ctx context.Context, profileID, characterID, categoryID *string) (bool, error)
}
