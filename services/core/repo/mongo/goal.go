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

func (r *GoalRepo) GetGoalsByCharacterID(ctx context.Context, characterID string) ([]entity.Goal, error) {
	return r.FindMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *GoalRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	return r.DeleteMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *GoalRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	return r.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
}

func (r *GoalRepo) Exist(ctx context.Context, characterID, goalID string) error {
	return r.Exists(ctx, bson.M{
		"_id":          mongodb.ToObjectID(goalID),
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *GoalRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}
