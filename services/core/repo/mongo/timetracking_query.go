package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *TimeTrackingRepo) FindByReferenceID(ctx context.Context, referenceID string) ([]entity.TimeTracking, error) {
	return r.FindMany(ctx, bson.M{
		"reference_id": mongodb.ToObjectID(referenceID),
	}, &options.FindOptions{
		Sort: bson.M{
			"timestamp": -1,
		},
	})
}

func (r *TimeTrackingRepo) FindByReferenceIDAndTimestamp(ctx context.Context, refID string, timestamp time.Time) (*entity.TimeTracking, error) {
	timetrack, err := r.FindOne(ctx, bson.M{
		"reference_id": mongodb.ToObjectID(refID),
		"timestamp":    bson.M{"$gte": utils.StartOfDay(timestamp), "$lt": utils.EndOfDay(timestamp)},
	})
	if err == mongo.ErrNoDocuments {
		return nil, errors.ErrMongoNotFound
	}

	return timetrack, err
}

func (r *TimeTrackingRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.TimeTracking, error) {
	return r.FindMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *TimeTrackingRepo) FindByCategoryID(ctx context.Context, categoryID string) ([]entity.TimeTracking, error) {
	return r.FindMany(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)})
}

func (r *TimeTrackingRepo) FindByTimeRange(ctx context.Context, start, end *time.Time) ([]entity.TimeTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if start != nil {
		filter["start_time"] = bson.M{"$gte": start}
	}
	if end != nil {
		filter["end_time"] = bson.M{"$lte": end}
	}

	return r.FindMany(ctx, filter)
}
