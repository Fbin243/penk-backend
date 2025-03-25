package mongorepo

import (
	"context"
	"log"
	"time"

	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"

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

	return &GoalRepo{mongodb.NewBaseRepo(
		goalCollection,
		&mongodb.Mapper[entity.Goal, mongomodel.Goal]{},
		true,
	)}
}

func (r *GoalRepo) GetGoalsByCharacterID(ctx context.Context, characterID string) ([]entity.Goal, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	goals := []entity.Goal{}
	err = cursor.All(ctx, &goals)

	return goals, err
}

func (r *GoalRepo) ValidateGoal(ctx context.Context, profileID, goalID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pineline := bson.A{
		bson.M{"$match": bson.M{"_id": mongodb.ToObjectID(goalID)}},
		bson.M{
			"$lookup": bson.M{
				"from":         mongodb.CharactersCollection,
				"localField":   "character_id",
				"foreignField": "_id",
				"as":           "character",
			},
		},
		bson.M{"$unwind": "$character"},
		bson.M{"$match": bson.M{"character.profile_id": mongodb.ToObjectID(profileID)}},
	}

	cursor, err := r.Aggregate(ctx, pineline)
	if err != nil || !cursor.Next(ctx) {
		return errors.ErrPermissionDenied
	}

	return nil
}

func (r *GoalRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	return err
}

func (r *GoalRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
	return err
}
