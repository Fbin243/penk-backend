package mongorepo

import (
	"context"
	"log"

	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoalRepo struct {
	*mongodb.BaseRepo[entity.Goal, mongomodel.Goal]
}

func NewGoalRepo(db *mongo.Database) *GoalRepo {
	goalCollection := db.Collection(mongodb.GoalsCollection)
	_, err := goalCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{Key: "character_id", Value: 1}},
	})
	if err != nil {
		log.Printf("failed to create indexes for %s collection\n", mongodb.GoalsCollection)
		return nil
	}

	return &GoalRepo{mongodb.NewBaseRepo[entity.Goal, mongomodel.Goal](
		goalCollection,
		true,
	)}
}
