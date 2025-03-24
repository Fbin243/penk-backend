package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TimeTrackingRepo) UnassignCategory(ctx context.Context, categoryID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.Collection.UpdateMany(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)}, bson.M{"$set": bson.M{"category_id": nil}})
	if err != nil {
		return err
	}

	return nil
}

func (r *TimeTrackingRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.Collection.DeleteMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	if err != nil {
		return err
	}

	return nil
}

func (r TimeTrackingRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.Collection.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
	if err != nil {
		return err
	}

	return nil
}
