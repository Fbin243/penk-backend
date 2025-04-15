package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/analytic/entity"

	rdb "tenkhours/pkg/db/redis"
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
func (biz *analyticBusiness) GetStatAnalytic(ctx context.Context, characterID string, startTime, endTime *time.Time, analyticSections []entity.AnalyticSection) (map[string]interface{}, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	capturedRecordFilter := entity.GetCapturedRecordFilter{
		CharacterID: authSession.CurrentCharacterID,
		EndTime:     time.Now(),
	}

	if startTime != nil {
		capturedRecordFilter.StartTime = utils.StartOfDay(*startTime)
	}

	if endTime != nil {
		capturedRecordFilter.EndTime = utils.StartOfDay(*endTime)
	}

	capturedRecords, err := biz.timetrackingRepo.AggregateDailyCapturedRecord(ctx, capturedRecordFilter)
	if err != nil {
		return nil, err
	}

	analyticsProcessor := &AnalyticsProcessor{
		AnalyticSections: analyticSections,
		CharacterID:      characterID,
		AnalyticResults:  make(map[string]any),
		CapturedRecords:  capturedRecords,
		StartTime:        capturedRecordFilter.StartTime,
		EndTime:          capturedRecordFilter.EndTime,
	}

	return analyticsProcessor.ProcessCapturedRecords(), nil
}
