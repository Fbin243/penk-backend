package business_test

import (
	"context"
	"os"
	"testing"
	"time"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/utils"
	"tenkhours/services/analytic/business"
	"tenkhours/services/analytic/entity"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockCapturedRecordRepo struct {
	mock.Mock
}

func (m *MockCapturedRecordRepo) GetCapturedRecords(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entity.CapturedRecord), args.Error(1)
}

func (m *MockCapturedRecordRepo) DeleteCapturedRecords(ctx context.Context, profileID string) error {
	args := m.Called(ctx, profileID)
	return args.Error(0)
}

type MockCoreClient struct {
	mock.Mock
}

func (m *MockCoreClient) CheckPermission(ctx context.Context, profileID, characterID, _ *string) (bool, error) {
	args := m.Called(ctx, profileID, characterID)
	return args.Bool(0), args.Error(1)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetCapturedRecords(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entity.CapturedRecord), args.Error(1)
}

type GetAnalyticResultsSuite struct {
	suite.Suite
	mockRepo       *MockCapturedRecordRepo
	mockCoreClient *MockCoreClient
	mockCache      *MockCache
	biz            business.IAnalyticBusiness
	profileID      string
	characterID    string
	startTime      time.Time
	endTime        time.Time
	ctx            context.Context
	filter         entity.GetCapturedRecordFilter
	sections       []entity.AnalyticSection
}

func (s *GetAnalyticResultsSuite) SetupSuite() {
	s.mockRepo = new(MockCapturedRecordRepo)
	s.mockCoreClient = new(MockCoreClient)
	s.mockCache = new(MockCache)
	s.biz = business.NewAnalyticBusiness(s.mockRepo, s.mockCoreClient, s.mockCache)
	s.profileID = "profile_id"
	s.characterID = "character_id_1"
	s.startTime = utils.ParseTime("2024-01-01T00:00:00.000Z")
	s.endTime = utils.ParseTime("2025-01-31T23:59:59.999Z")
	s.ctx = context.WithValue(context.Background(), auth.AuthSessionKey, rdb.AuthSession{
		ProfileID: s.profileID,
	})
	s.filter = entity.GetCapturedRecordFilter{
		ProfileID:   s.profileID,
		CharacterID: &s.characterID,
		StartTime:   utils.ResetTimeToBeginningOfDay(s.startTime),
		EndTime:     utils.ResetTimeToBeginningOfDay(s.endTime),
	}
	s.sections = []entity.AnalyticSection{
		entity.AnalyticSectionDistribution,
		entity.AnalyticSectionFrequency,
		entity.AnalyticSectionOverall,
		entity.AnalyticSectionTimeline,
	}
}

func (s *GetAnalyticResultsSuite) TestGetAnalyticResultsForProfile() {
	s.filter.CharacterID = nil
	s.mockRepo.On("GetCapturedRecords", s.ctx, s.filter).Return([]entity.CapturedRecord{capturedRecords[0]}, nil)
	s.mockCache.On("GetCapturedRecords", s.ctx, s.filter).Return([]entity.CapturedRecord{capturedRecords[1]}, nil)

	result, err := s.biz.GetAnalyticResults(s.ctx, nil, &s.startTime, &s.endTime, s.sections)

	s.NoError(err)
	s.NotNil(result)
	// log.Printf("result: %v", utils.PrettyJSON(result))
	analyticResultJSON, _ := os.ReadFile("profile_analytic_result.json")
	s.JSONEq(string(analyticResultJSON), utils.PrettyJSON(result))

	s.mockRepo.AssertExpectations(s.T())
	s.mockCache.AssertExpectations(s.T())
}

func (s *GetAnalyticResultsSuite) TestGetAnalyticResultsForCharacter() {
	s.mockCoreClient.On("CheckPermission", s.ctx, &s.profileID, &s.characterID).Return(true, nil)
	s.mockRepo.On("GetCapturedRecords", s.ctx, s.filter).Return([]entity.CapturedRecord{capturedRecords[0]}, nil)
	s.mockCache.On("GetCapturedRecords", s.ctx, s.filter).Return([]entity.CapturedRecord{}, nil)

	result, err := s.biz.GetAnalyticResults(s.ctx, &s.characterID, &s.startTime, &s.endTime, s.sections)

	s.NoError(err)
	s.NotNil(result)
	// log.Printf("result: %v", utils.PrettyJSON(result))
	analyticResultJSON, _ := os.ReadFile("character_analytic_result.json")
	s.JSONEq(string(analyticResultJSON), utils.PrettyJSON(result))

	s.mockCoreClient.AssertExpectations(s.T())
	s.mockRepo.AssertExpectations(s.T())
	s.mockCache.AssertExpectations(s.T())
}

func TestGetAnalyticResultsSuite(t *testing.T) {
	suite.Run(t, new(GetAnalyticResultsSuite))
}
