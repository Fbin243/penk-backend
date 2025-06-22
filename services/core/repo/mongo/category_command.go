package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *CategoryRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	return r.DeleteMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *CategoryRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	return r.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
}
