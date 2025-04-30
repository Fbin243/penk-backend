package business

import (
	"context"

	"tenkhours/services/analytic/entity"
)

// Business
type IAnalyticBusiness interface {
	GetStatAnalytic(ctx context.Context, filter *entity.StatAnalyticFilter) (map[string]any, error)
}

// RPC client
type ICoreClient interface {
	CheckPermission(ctx context.Context, profileID, characterID string) (bool, error)
}

type ITimeTrackingRepo interface {
	AggregateDailyCapturedRecord(ctx context.Context, filter entity.StatAnalyticFilter) ([]entity.CapturedRecord, error)
}
