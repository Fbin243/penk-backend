package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/analytic/entity"

	rdb "tenkhours/pkg/db/redis"

	"github.com/samber/lo"
)

type analyticBusiness struct {
	coreClient       ICoreClient
	timetrackingRepo ITimeTrackingRepo
}

func NewAnalyticBusiness(coreClient ICoreClient, timetrackingRepo ITimeTrackingRepo) *analyticBusiness {
	return &analyticBusiness{
		coreClient:       coreClient,
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
		// ProfileID:   authSession.ProfileID,
		CharacterID: characterID,
		EndTime:     time.Now(),
	}

	authorized, err := biz.coreClient.CheckPermission(ctx, lo.ToPtr(authSession.ProfileID), lo.ToPtr(characterID), nil)
	if !authorized || err != nil {
		return nil, errors.ErrPermissionDenied
	}

	if startTime != nil {
		capturedRecordFilter.StartTime = utils.ResetTimeToBeginningOfDay(*startTime)
	}

	if endTime != nil {
		capturedRecordFilter.EndTime = utils.ResetTimeToBeginningOfDay(*endTime)
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
