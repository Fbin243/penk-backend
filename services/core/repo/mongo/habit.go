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

type HabitRepo struct {
	*mongodb.BaseRepo[entity.Habit, mongomodel.Habit]
}

func NewHabitRepo(db *mongo.Database) *HabitRepo {
	habitCollection := db.Collection(mongodb.HabitsCollection)
	_, err := habitCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "character_id", Value: 1}}},
		{Keys: bson.D{{Key: "category_id", Value: 1}}},
	})
	if err != nil {
		log.Printf("failed to create indexes for %s collection\n", mongodb.HabitsCollection)
		return nil
	}

	return &HabitRepo{mongodb.NewBaseRepo[entity.Habit, mongomodel.Habit](
		habitCollection,
		true,
	)}
}
