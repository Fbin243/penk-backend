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

	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTimeTrackingRepo struct {
	mock.Mock
}

func (m *MockTimeTrackingRepo) AggregateDailyCapturedRecord(ctx context.Context, filter entity.StatAnalyticFilter) ([]entity.CapturedRecord, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entity.CapturedRecord), args.Error(1)
}

func (m *MockTimeTrackingRepo) DeleteCapturedRecords(ctx context.Context, profileID string) error {
	args := m.Called(ctx, profileID)
	return args.Error(0)
}

type GetAnalyticResultsSuite struct {
	suite.Suite
	mockRepo    *MockTimeTrackingRepo
	biz         business.IAnalyticBusiness
	profileID   string
	characterID string
	startTime   time.Time
	endTime     time.Time
	ctx         context.Context
	filter      entity.StatAnalyticFilter
}

func (s *GetAnalyticResultsSuite) SetupSuite() {
	s.mockRepo = new(MockTimeTrackingRepo)
	s.biz = business.NewAnalyticBusiness(s.mockRepo)
	s.profileID = "profile_id"
	s.characterID = "character_id_1"
	s.startTime = utils.ParseTime("2025-01-01T00:00:00.000Z")
	s.endTime = utils.ParseTime("2025-12-31T23:59:59.999Z")
	s.ctx = context.WithValue(context.Background(), auth.AuthSessionKey, rdb.AuthSession{
		ProfileID: s.profileID,
	})
	s.filter = entity.StatAnalyticFilter{
		CharacterID: s.characterID,
		StartTime:   lo.ToPtr(utils.StartOfDay(s.startTime)),
		EndTime:     lo.ToPtr(utils.StartOfDay(s.endTime)),
		AnalyticSections: []entity.AnalyticSection{
			entity.AnalyticSectionDistribution,
			entity.AnalyticSectionFrequency,
			entity.AnalyticSectionOverall,
			entity.AnalyticSectionTimeline,
		},
	}
}

func (s *GetAnalyticResultsSuite) TestGetAnalyticResultsForCharacter() {
	s.mockRepo.On("AggregateDailyCapturedRecord", s.ctx, s.filter).Return(capturedRecords, nil)

	result, err := s.biz.GetStatAnalytic(s.ctx, &s.filter)

	s.NoError(err)
	s.NotNil(result)
	analyticResultJSON, _ := os.ReadFile("character_analytic_result.json")
	s.JSONEq(string(analyticResultJSON), utils.PrettyJSON(result))

	s.mockRepo.AssertExpectations(s.T())
}

func TestGetAnalyticResultsSuite(t *testing.T) {
	suite.Run(t, new(GetAnalyticResultsSuite))
}
