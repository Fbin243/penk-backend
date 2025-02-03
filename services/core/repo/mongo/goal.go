package mongorepo

import (
	"context"
	"time"

	"tenkhours/services/core/entity"

	mongodb "tenkhours/pkg/db/mongo"

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

func (r *GoalRepo) GetGoalsByCharacterID(ctx context.Context, characterID string, status *entity.GoalStatusFilter) ([]entity.Goal, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"character_id": characterID}
	if status != nil {
		if status.FinishStatus != nil {
			filter["status"] = status.FinishStatus
		}
		if status.ExpireStatus != nil {
			switch *status.ExpireStatus {
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

func (r *GoalRepo) UpdateOneMetricInGoals(ctx context.Context, metric entity.CustomMetric, goalIDs []string) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": goalIDs}, "target": bson.M{"$elemMatch": bson.M{"_id": metric.ID}}}
	update := bson.M{
		"$set": bson.M{
			"target.$.name":        metric.Name,
			"target.$.description": metric.Description,
			"target.$.style":       metric.Style,
		},
	}

	return r.UpdateMany(ctx, filter, update)
}

func (r *GoalRepo) RemoveOnePropertyFromGoals(ctx context.Context, metricID, propertyID string, goalIDs []string) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": bson.M{"$in": goalIDs}, "target": bson.M{"$elemMatch": bson.M{"_id": metricID}}}
	update := bson.M{"$pull": bson.M{"target.$.properties": bson.M{"_id": propertyID}}}

	return r.UpdateMany(ctx, filter, update)
}

func (r *GoalRepo) UpdateStatusOfGoals(ctx context.Context, goalIDs []string, status entity.GoalFinishStatus) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": goalIDs}}
	update := bson.M{"$set": bson.M{"status": status}}

	return r.UpdateMany(ctx, filter, update)
}
