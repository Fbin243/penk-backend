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

type TaskSessionRepo struct {
	*mongodb.BaseRepo[entity.TaskSession, mongomodel.TaskSession]
}

func NewTaskSessionRepo(db *mongo.Database) *TaskSessionRepo {
	taskSessionCollection := db.Collection(mongodb.TaskSessionsCollection)
	_, err := taskSessionCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "task_id", Value: 1}}},
	})
	if err != nil {
		log.Printf("failed to create indexes for %s collection\n", mongodb.TaskSessionsCollection)
		return nil
	}

	return &TaskSessionRepo{mongodb.NewBaseRepo[entity.TaskSession, mongomodel.TaskSession](
		taskSessionCollection,
		false,
	)}
}
