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
	capturedRecordRepo ICapturedRecordRepo
	coreClient         ICoreClient
	cache              ICache
}

func NewAnalyticBusiness(capturedRepo ICapturedRecordRepo, coreClient ICoreClient, cache ICache) *analyticBusiness {
	return &analyticBusiness{
		capturedRecordRepo: capturedRepo,
		coreClient:         coreClient,
		cache:              cache,
	}
}

// Get analytic results from the captured records from the database
func (biz *analyticBusiness) GetAnalyticResults(ctx context.Context, characterID *string, startTime, endTime *time.Time, analyticSections []entity.AnalyticSection) (map[string]interface{}, error) {
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
		authorized, err := biz.coreClient.CheckPermission(ctx, lo.ToPtr(authSession.ProfileID), characterID, nil)
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

	// Get captured records from db
	capturedRecords, err := biz.capturedRecordRepo.GetCapturedRecords(ctx, capturedRecordFilter)
	if err != nil {
		return nil, err
	}

	// Get captured records from cache
	cacheCapturedRecords, err := biz.cache.GetCapturedRecords(ctx, capturedRecordFilter)
	if err != nil && err != errors.ErrRedisNotFound {
		return nil, err
	}

	analyticsProcessor := &AnalyticsProcessor{
		AnalyticSections: analyticSections,
		ProfileID:        authSession.ProfileID,
		CharacterID:      characterID,
		CapturedRecords:  append(capturedRecords, cacheCapturedRecords...),
		AnalyticResults:  make(map[string]interface{}),
		StartTime:        capturedRecordFilter.StartTime,
		EndTime:          capturedRecordFilter.EndTime,
	}

	return analyticsProcessor.ProcessCapturedRecords(), nil
}

func (biz *analyticBusiness) DeleteCapturedRecords(ctx context.Context, profileID string) error {
	return biz.capturedRecordRepo.DeleteCapturedRecords(ctx, profileID)
}
