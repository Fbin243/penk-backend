package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TimeTrackingRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.TimeTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	timeTrackings := []entity.TimeTracking{}
	err = cursor.All(ctx, &timeTrackings)
	if err != nil {
		return nil, err
	}

	return timeTrackings, nil
}

func (r *TimeTrackingRepo) FindByCategoryID(ctx context.Context, categoryID string) ([]entity.TimeTracking, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	timeTrackings := []entity.TimeTracking{}
	err = cursor.All(ctx, &timeTrackings)
	if err != nil {
		return nil, err
	}

	return timeTrackings, nil
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

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	timeTrackings := []entity.TimeTracking{}
	err = cursor.All(ctx, &timeTrackings)
	if err != nil {
		return nil, err
	}

	return timeTrackings, nil
}
