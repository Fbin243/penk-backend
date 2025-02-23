package mongorepo_test

import (
	"context"
	"testing"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/services/analytic/entity"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	profileID   = mongodb.GenObjectID()
	characterID = mongodb.GenObjectID()
	categoryID  = mongodb.GenObjectID()
	now         = utils.Now()
)

func NewCapturedRecord() *entity.CapturedRecord {
	return &entity.CapturedRecord{
		ID:               mongodb.GenObjectID(),
		Timestamp:        now,
		TotalFocusedTime: 100,
		Categories: []entity.CapturedRecordCategory{
			{
				ID:   categoryID,
				Time: 100,
			},
		},
		TimeTrackings: []entity.CapturedRecordTimeTracking{
			{
				CategoryID: lo.ToPtr(categoryID),
				Time:       100,
				StartTime:  now,
				EndTime:    now,
			},
		},
		Metadata: entity.CapturedRecordMetadata{
			CharacterID: characterID,
			ProfileID:   profileID,
		},
	}
}

func TestCreateCapturedRecord(t *testing.T) {
	capturedRecord := NewCapturedRecord()
	createdCapturedRecord, err := capturedRecordRepo.CreateCapturedRecord(context.Background(), capturedRecord)
	defer cleanUp(t, capturedRecord.ID)
	assert.Nil(t, err)
	assert.Equal(t, capturedRecord, createdCapturedRecord)
}

func TestGetCapturedRecords(t *testing.T) {
	suite.Run(t, new(GetCapturedRecordsSuite))
}

func TestDeleteCapturedRecords(t *testing.T) {
	capturedRecord := NewCapturedRecord()
	_, err := capturedRecordRepo.CreateCapturedRecord(context.Background(), capturedRecord)
	assert.Nil(t, err)

	err = capturedRecordRepo.DeleteCapturedRecords(context.Background(), profileID)
	assert.Nil(t, err)

	filter := entity.GetCapturedRecordFilter{
		ProfileID: profileID,
		StartTime: now.Add(-time.Hour),
		EndTime:   now.Add(time.Hour),
	}
	capturedRecords, err := capturedRecordRepo.GetCapturedRecords(context.Background(), filter)
	assert.Nil(t, err)
	assert.Empty(t, capturedRecords)
}

type GetCapturedRecordsSuite struct {
	suite.Suite
	capturedRecords []entity.CapturedRecord
}

func (s *GetCapturedRecordsSuite) SetupSuite() {
	for i := 0; i < 2; i++ {
		capturedRecord := NewCapturedRecord()
		s.capturedRecords = append(s.capturedRecords, *capturedRecord)
	}

	s.capturedRecords[1].Metadata.CharacterID = mongodb.GenObjectID()

	for _, capturedRecord := range s.capturedRecords {
		_, err := capturedRecordRepo.CreateCapturedRecord(context.Background(), &capturedRecord)
		s.Nil(err)
	}
}

func (s *GetCapturedRecordsSuite) TearDownSuite() {
	for _, capturedRecord := range s.capturedRecords {
		cleanUp(s.T(), capturedRecord.ID)
	}
}

func (s *GetCapturedRecordsSuite) TestFilterWithProfile() {
	filter := entity.GetCapturedRecordFilter{
		ProfileID: profileID,
		StartTime: now.Add(-time.Hour),
		EndTime:   now.Add(time.Hour),
	}

	capturedRecords, err := capturedRecordRepo.GetCapturedRecords(context.Background(), filter)
	s.NotEmpty(capturedRecords)
	s.Nil(err)
	s.Len(capturedRecords, 2)
	s.Equal(s.capturedRecords, capturedRecords)
}

func (s *GetCapturedRecordsSuite) TestFilterWithCharacter() {
	filter := entity.GetCapturedRecordFilter{
		ProfileID:   profileID,
		CharacterID: &characterID,
		StartTime:   now.Add(-time.Hour),
		EndTime:     now.Add(time.Hour),
	}

	capturedRecords, err := capturedRecordRepo.GetCapturedRecords(context.Background(), filter)
	s.NotEmpty(capturedRecords)
	s.Nil(err)
	s.Len(capturedRecords, 1)
	s.Equal(s.capturedRecords[0], capturedRecords[0])
}

func (s *GetCapturedRecordsSuite) TestFilterWithTimeRange() {
	filter := entity.GetCapturedRecordFilter{
		ProfileID: profileID,
		StartTime: now.Add(-time.Hour),
		EndTime:   now.Add(-time.Minute),
	}

	capturedRecords, err := capturedRecordRepo.GetCapturedRecords(context.Background(), filter)
	s.Empty(capturedRecords)
	s.Nil(err)
}

func cleanUp(t *testing.T, id string) {
	_, err := capturedRecordRepo.DeleteOne(context.Background(), bson.M{"_id": mongodb.ToObjectID(id)})
	assert.Nil(t, err)
}
