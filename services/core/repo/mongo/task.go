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

type TaskRepo struct {
	*mongodb.BaseRepo[entity.Task, mongomodel.Task]
}

func NewTaskRepo(db *mongo.Database) *TaskRepo {
	taskCollection := db.Collection(mongodb.TasksCollection)
	_, err := taskCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "character_id", Value: 1}}},
	})
	if err != nil {
		log.Printf("failed to create indexes for %s collection\n", mongodb.TasksCollection)
		return nil
	}

	return &TaskRepo{mongodb.NewBaseRepo[entity.Task, mongomodel.Task](
		taskCollection,
		true,
	)}
}
