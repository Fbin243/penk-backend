package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/services/analytic/entity"
)

type analyticBusiness struct {
	timetrackingRepo ITimeTrackingRepo
}

func NewAnalyticBusiness(timetrackingRepo ITimeTrackingRepo) *analyticBusiness {
	return &analyticBusiness{
		timetrackingRepo: timetrackingRepo,
	}
}

// Get analytic results from the captured records from the database
func (biz *analyticBusiness) GetStatAnalytic(ctx context.Context, filter *entity.StatAnalyticFilter) (map[string]interface{}, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = &entity.StatAnalyticFilter{}
	}
	filter.CharacterID = authSession.CurrentCharacterID

	capturedRecords, err := biz.timetrackingRepo.AggregateDailyCapturedRecord(ctx, *filter)
	if err != nil {
		return nil, err
	}

	analyticsProcessor := &AnalyticsProcessor{
		AnalyticSections: filter.AnalyticSections,
		CharacterID:      authSession.CurrentCharacterID,
		AnalyticResults:  make(map[string]any),
		CapturedRecords:  capturedRecords,
		StartTime:        filter.StartTime,
		EndTime:          filter.EndTime,
	}

	return analyticsProcessor.ProcessCapturedRecords()
}
