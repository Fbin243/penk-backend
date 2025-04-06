package business

import (
	"context"
	"time"

	"tenkhours/services/analytic/entity"
)

// Business
type IAnalyticBusiness interface {
	GetStatAnalytic(ctx context.Context, characterID string, startTime, endTime *time.Time, analyticSections []entity.AnalyticSection) (map[string]any, error)
}

// RPC client
type ICoreClient interface {
	CheckPermission(ctx context.Context, profileID, characterID string) (bool, error)
}

type ITimeTrackingRepo interface {
	AggregateDailyCapturedRecord(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error)
}
