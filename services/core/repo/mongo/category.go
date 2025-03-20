package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
	mongorepo "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepo struct {
	*mongodb.BaseRepo[entity.Category, mongorepo.Category]
}

func NewCategoryRepo(db *mongo.Database) *CategoryRepo {
	return &CategoryRepo{mongodb.NewBaseRepo(
		db.Collection(mongodb.CategoriesCollection),
		&mongodb.Mapper[entity.Category, mongorepo.Category]{},
	)}
}

func (r *CategoryRepo) ValidateCategory(ctx context.Context, characterID, categoryID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID), "_id": mongodb.ToObjectID(categoryID)})
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
