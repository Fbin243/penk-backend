package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *HabitLogRepo) FindByHabitID(ctx context.Context, habitID string, filter *entity.HabitLogFilter, orderBy *entity.HabitLogOrderBy, limit, offset *int) ([]entity.HabitLog, error) {
	mongoFilter := bson.M{
		"habit_id": mongodb.ToObjectID(habitID),
	}

	if filter != nil {
		mongoFilter["timestamp"] = bson.M{}

		if filter.StartTime != nil {
			mongoFilter["timestamp"].(bson.M)["$gte"] = filter.StartTime.Format(time.DateOnly)
		}
		if filter.EndTime != nil {
			mongoFilter["timestamp"].(bson.M)["$lte"] = filter.EndTime.Format(time.DateOnly)
		}
	}

	opts := options.Find()
	if orderBy != nil {
		if orderBy.Timestamp != nil {
			opts.SetSort(bson.M{"timestamp": orderBy.Timestamp.ToInt()})
		}
	}
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.FindMany(ctx, mongoFilter, opts)
}

func (r *HabitLogRepo) FindByHabitIDs(ctx context.Context, habitIDs []string, filter *entity.HabitLogFilter, orderBy *entity.HabitLogOrderBy, limit, offset *int) ([]entity.HabitLog, error) {
	mongoFilter := bson.M{
		"habit_id": bson.M{
			"$in": mongodb.ToObjectIDs(habitIDs),
		},
	}

	if filter != nil {
		mongoFilter["timestamp"] = bson.M{}

		if filter.StartTime != nil {
			mongoFilter["timestamp"].(bson.M)["$gte"] = filter.StartTime.Format(time.DateOnly)
		}
		if filter.EndTime != nil {
			mongoFilter["timestamp"].(bson.M)["$lte"] = filter.EndTime.Format(time.DateOnly)
		}
	}

	opts := options.Find()
	if orderBy != nil {
		if orderBy.Timestamp != nil {
			opts.SetSort(bson.M{"timestamp": orderBy.Timestamp.ToInt()})
		}
	}
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.FindMany(ctx, mongoFilter, opts)
}
