package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepo struct {
	*mongodb.BaseRepo[entity.Category, mongomodel.Category]
}

func NewCategoryRepo(db *mongo.Database) *CategoryRepo {
	return &CategoryRepo{mongodb.NewBaseRepo(
		db.Collection(mongodb.CategoriesCollection),
		&mongodb.Mapper[entity.Category, mongomodel.Category]{},
		true,
	)}
}

func (r *CategoryRepo) CountByCharacterID(ctx context.Context, characterID string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *CategoryRepo) Exist(ctx context.Context, characterID, categoryID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.ErrMongoNotFound
	}

	return nil
}

func (r *CategoryRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	categories := []entity.Category{}
	err = cursor.All(ctx, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	return err
}

func (r *CategoryRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
	return err
}
