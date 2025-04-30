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
