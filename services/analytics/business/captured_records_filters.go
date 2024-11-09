package business

import (
	"sort"
	"time"

	"tenkhours/services/analytics/graph/model"
	"tenkhours/services/analytics/repo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilterMethod int

const (
	FilterMethodServerSide FilterMethod = iota
	FilterMethodClientSide
)

type FilterType int

const (
	FilterTypeUser FilterType = iota
	FilterTypeCharacter
)

type CapturedRecordsFilter struct {
	FilterMethod         FilterMethod
	FilterType           FilterType
	CharacterID          primitive.ObjectID
	ProfileID            primitive.ObjectID
	StartTime            time.Time
	EndTime              time.Time
	CapturedRecordsRepo  *repo.CapturedRecordsRepo
	CapturedRecordLocals []model.CapturedRecord
}

func (c *CapturedRecordsFilter) Filter() ([]model.CapturedRecord, error) {
	switch c.FilterMethod {
	case FilterMethodServerSide:
		matchStage := bson.D{}
		if c.FilterType == FilterTypeCharacter {
			matchStage = append(matchStage, bson.E{Key: "metadata.character_id", Value: c.CharacterID})
		} else {
			matchStage = append(matchStage, bson.E{Key: "metadata.profile_id", Value: c.ProfileID})
		}

		matchStage = append(matchStage, bson.E{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: c.StartTime}}})
		matchStage = append(matchStage, bson.E{Key: "timestamp", Value: bson.D{{Key: "$lte", Value: c.EndTime}}})

		pipeline := mongo.Pipeline{}
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchStage}})
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "timestamp", Value: 1}}}})

		return c.CapturedRecordsRepo.GetCapturedRecords(pipeline)

	case FilterMethodClientSide:
		capturedRecordLocals := c.CapturedRecordLocals
		if c.FilterType == FilterTypeCharacter {
			capturedRecordLocals = lo.Filter(capturedRecordLocals, func(record model.CapturedRecord, _ int) bool {
				return record.Metadata.CharacterID == c.CharacterID
			})
		}

		capturedRecordLocals = lo.Filter(capturedRecordLocals, func(record model.CapturedRecord, _ int) bool {
			return record.Timestamp.After(c.StartTime.Add(-24*time.Hour)) && record.Timestamp.Before(c.EndTime.Add(24*time.Hour))
		})

		sort.Slice(capturedRecordLocals, func(i, j int) bool {
			return capturedRecordLocals[i].Timestamp.Before(capturedRecordLocals[j].Timestamp)
		})

		return capturedRecordLocals, nil
	}

	return nil, nil
}
