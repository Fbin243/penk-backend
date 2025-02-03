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
	CreateCapturedRecord(ctx context.Context, capturedRecord *entity.CapturedRecord) error
	GetCapturedRecords(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error)
	DeleteCapturedRecords(ctx context.Context, profileID string) error
}

// RPC client
type ICoreClient interface {
	CheckPermission(ctx context.Context, profileID, characterID string, metricID *string) (bool, error)
}
