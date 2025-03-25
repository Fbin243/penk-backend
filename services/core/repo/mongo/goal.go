package mongorepo

import (
	"context"
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
	return &GoalRepo{mongodb.NewBaseRepo(
		db.Collection(mongodb.GoalsCollection),
		&mongodb.Mapper[entity.Goal, mongomodel.Goal]{},
		true,
	)}
}

func (r *GoalRepo) GetGoalsByCharacterID(ctx context.Context, characterID string, status *entity.GoalStatus) ([]entity.Goal, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.SyncGoalStatus(ctx, characterID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"character_id": mongodb.ToObjectID(characterID)}
	if status != nil {
		filter["status"] = status
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

func (r *GoalRepo) UpdateStatusOfGoals(ctx context.Context, goalIDs []string, status entity.GoalStatus) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": mongodb.ToObjectIDs(goalIDs)}}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.UpdateMany(ctx, filter, update)

	return err
}

func (r *GoalRepo) SyncGoalStatus(ctx context.Context, characterID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	operations := []mongo.WriteModel{
		// Planned --> In Progress
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{
				"character_id": mongodb.ToObjectID(characterID),
				"status":       entity.GoalStatusPlanned,
				"start_time":   bson.M{"$lte": time.Now()},
			}).
			SetUpdate(bson.M{
				"$set": bson.M{
					"status": entity.GoalStatusInProgress,
				},
			}),

		// In Progress --> Overdue
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{
				"character_id": mongodb.ToObjectID(characterID),
				"status":       entity.GoalStatusInProgress,
				"end_time":     bson.M{"$lt": time.Now()},
			}).
			SetUpdate(bson.M{
				"$set": bson.M{
					"status": entity.GoalStatusOverdue,
				},
			}),
	}

	_, err := r.BulkWrite(ctx, operations)

	return err
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
