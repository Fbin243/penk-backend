package repo

import (
	"context"
	"time"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoalsRepo struct {
	*db.BaseRepo[Goal]
}

func NewGoalsRepo(mongo *mongo.Database) *GoalsRepo {
	return &GoalsRepo{db.NewBaseRepo[Goal](mongo.Collection(db.GoalsCollection))}
}

func (r *GoalsRepo) GetGoalsByCharacterID(characterID primitive.ObjectID, status *GoalStatusFilter) ([]Goal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"character_id": characterID}
	if status != nil {
		if status.FinishStatus != nil {
			filter["status"] = status.FinishStatus
		}
		if status.ExpireStatus != nil {
			switch *status.ExpireStatus {
			case GoalExpireStatusExpired:
				filter["end_date"] = bson.M{"$lte": time.Now()}
			case GoalExpireStatusUnexpired:
				filter["end_date"] = bson.M{"$gt": time.Now()}
			}
		}
	}

	goals := []Goal{}
	cursor, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &goals)

	return goals, err
}

func (r *GoalsRepo) RemoveOneMetricFromGoals(metricID primitive.ObjectID, goalIDs []primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.UpdateMany(ctx,
		bson.M{"_id": bson.M{"$in": goalIDs}},
		bson.M{"$pull": bson.M{"target": bson.M{"_id": metricID}}})

	return err
}
