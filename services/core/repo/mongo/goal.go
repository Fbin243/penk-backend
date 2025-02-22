package mongorepo

import (
	"context"
	"time"

	"tenkhours/services/core/entity"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoalRepo struct {
	*mongodb.BaseRepo[entity.Goal, Goal]
}

func NewGoalRepo(db *mongo.Database) *GoalRepo {
	return &GoalRepo{mongodb.NewBaseRepo(
		db.Collection(mongodb.GoalsCollection),
		&mongodb.Mapper[entity.Goal, Goal]{},
	)}
}

func (r *GoalRepo) GetGoalsByCharacterID(ctx context.Context, characterID string, statusFilter *entity.GoalStatusFilter) ([]entity.Goal, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"character_id": mongodb.ToObjectID(characterID)}
	if statusFilter != nil {
		if statusFilter.FinishStatus != nil {
			filter["status"] = statusFilter.FinishStatus
		}
		if statusFilter.ExpireStatus != nil {
			switch *statusFilter.ExpireStatus {
			case entity.GoalExpireStatusExpired:
				filter["end_date"] = bson.M{"$lte": time.Now()}
			case entity.GoalExpireStatusUnexpired:
				filter["end_date"] = bson.M{"$gt": time.Now()}
			}
		}
	}

	goals := []entity.Goal{}
	cursor, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

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

func (r *GoalRepo) UpdateStatusOfGoals(ctx context.Context, goalIDs []string, status entity.GoalFinishStatus) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": mongodb.ToObjectIDs(goalIDs)}}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.UpdateMany(ctx, filter, update)

	return err
}
