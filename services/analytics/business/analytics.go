package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/analytics/graph/model"
	analyticsRepo "tenkhours/services/analytics/repo"
	coreRepo "tenkhours/services/core/repo"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalyticsBusiness struct {
	CharactersRepo      *coreRepo.CharactersRepo
	ProfilesRepo        *coreRepo.ProfilesRepo
	CapturedRecordsRepo *analyticsRepo.CapturedRecordsRepo
	RedisClient         *redis.Client
}

func NewAnalyticsBusiness(charactersRepo *coreRepo.CharactersRepo, profilesRepo *coreRepo.ProfilesRepo, capturedRepo *analyticsRepo.CapturedRecordsRepo, redisClient *redis.Client) *AnalyticsBusiness {
	return &AnalyticsBusiness{
		CharactersRepo:      charactersRepo,
		ProfilesRepo:        profilesRepo,
		CapturedRecordsRepo: capturedRepo,
		RedisClient:         redisClient,
	}
}

// Get analytic results from the captured records from the database
func (biz *AnalyticsBusiness) GetAnalyticResults(ctx context.Context, characterID *primitive.ObjectID, startTime *time.Time, endTime *time.Time, analyticSections []model.AnalyticSection) (map[string]interface{}, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	capturedRecordsFilter := CapturedRecordsFilter{
		FilterMethod:        FilterMethodServerSide,
		FilterType:          FilterTypeUser,
		ProfileID:           authSession.ProfileID,
		EndTime:             time.Now(),
		CapturedRecordsRepo: biz.CapturedRecordsRepo,
		RedisClient:         biz.RedisClient,
	}

	if characterID != nil {
		character, err := biz.CharactersRepo.FindByID(*characterID)
		if err != nil {
			return nil, err
		}

		if character.ProfileID != authSession.ProfileID {
			return nil, errors.PermissionDenied()
		}

		capturedRecordsFilter.FilterType = FilterTypeCharacter
		capturedRecordsFilter.CharacterID = *characterID
	}

	if startTime != nil {
		capturedRecordsFilter.StartTime = utils.ResetTimeToBeginningOfDay(*startTime)
	}

	if endTime != nil {
		capturedRecordsFilter.EndTime = utils.ResetTimeToBeginningOfDay(*endTime)
	}

	capturedRecords, err := capturedRecordsFilter.Filter()
	if err != nil {
		return nil, err
	}

	analyticsProcessor := &AnalyticsProcessor{
		AnalyticSections: analyticSections,
		CapturedRecords:  capturedRecords,
		AnalyticResults:  make(map[string]interface{}),
		FilterType:       capturedRecordsFilter.FilterType,
		StartTime:        capturedRecordsFilter.StartTime,
		EndTime:          capturedRecordsFilter.EndTime,
	}

	return analyticsProcessor.ProcessCapturedRecords(), nil
}
