package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

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
