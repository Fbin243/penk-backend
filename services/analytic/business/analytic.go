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

type AnalyticBusiness struct {
	CapturedRecordRepo ICapturedRecordRepo
	CoreClient         ICoreClient
}

func NewAnalyticBusiness(capturedRepo ICapturedRecordRepo, coreClient ICoreClient) *AnalyticBusiness {
	return &AnalyticBusiness{
		CapturedRecordRepo: capturedRepo,
		CoreClient:         coreClient,
	}
}

// Get analytic results from the captured records from the database
func (biz *AnalyticBusiness) GetAnalyticResults(ctx context.Context, characterID *string, startTime, endTime *time.Time, analyticSections []entity.AnalyticSection) (map[string]interface{}, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	capturedRecordFilter := entity.GetCapturedRecordFilter{
		ProfileID:   authSession.ProfileID,
		CharacterID: characterID,
		EndTime:     time.Now(),
	}

	if characterID != nil {
		authorized, err := biz.CoreClient.CheckPermission(ctx, lo.ToPtr(authSession.ProfileID), characterID, nil)
		if !authorized || err != nil {
			return nil, errors.ErrPermissionDenied
		}
	}

	if startTime != nil {
		capturedRecordFilter.StartTime = utils.ResetTimeToBeginningOfDay(*startTime)
	}

	if endTime != nil {
		capturedRecordFilter.EndTime = utils.ResetTimeToBeginningOfDay(*endTime)
	}

	// Get captured records from db and cache
	capturedRecords, err := biz.CapturedRecordRepo.GetCapturedRecords(ctx, capturedRecordFilter)
	if err != nil {
		return nil, err
	}

	analyticsProcessor := &AnalyticsProcessor{
		AnalyticSections: analyticSections,
		CapturedRecords:  capturedRecords,
		AnalyticResults:  make(map[string]interface{}),
		StartTime:        capturedRecordFilter.StartTime,
		EndTime:          capturedRecordFilter.EndTime,
	}

	return analyticsProcessor.ProcessCapturedRecords(), nil
}

func (biz *AnalyticBusiness) DeleteCapturedRecords(ctx context.Context, profileID string) error {
	return biz.CapturedRecordRepo.DeleteCapturedRecords(ctx, profileID)
}
