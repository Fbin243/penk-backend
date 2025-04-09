package mongorepo

import (
	"context"
	"log"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepo struct {
	*mongodb.BaseRepo[entity.Category, mongomodel.Category]
}

func NewCategoryRepo(db *mongo.Database) *CategoryRepo {
	cateColl := db.Collection(mongodb.CategoriesCollection)
	_, err := cateColl.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "character_id", Value: 1}}},
	})
	if err != nil {
		log.Printf("failed to create indexes for %s collection\n", mongodb.CategoriesCollection)
		return nil
	}

	return &CategoryRepo{mongodb.NewBaseRepo[entity.Category, mongomodel.Category](cateColl, true)}
}

func (r *CategoryRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *CategoryRepo) Exist(ctx context.Context, characterID, categoryID string) error {
	return r.Exists(ctx, bson.M{
		"_id":          mongodb.ToObjectID(categoryID),
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *CategoryRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.Category, error) {
	return r.FindMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *CategoryRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	return r.DeleteMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *CategoryRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	return r.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
}
