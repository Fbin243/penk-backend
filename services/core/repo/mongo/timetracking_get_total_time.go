package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TimeTrackingRepo) GetTotalTimeByCategoryID(ctx context.Context, categoryID string) (int, error) {
	return r.getTotalTime(ctx, bson.M{
		"category_id": mongodb.ToObjectID(categoryID),
	})
}

func (r *TimeTrackingRepo) GetTotalTimeByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.getTotalTime(ctx, bson.M{
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *TimeTrackingRepo) GetTotalTimeOfUnassigned(ctx context.Context, characterID string) (int, error) {
	return r.getTotalTime(ctx, bson.M{
		"character_id": mongodb.ToObjectID(characterID),
		"category_id":  nil,
	})
}

func (r *TimeTrackingRepo) getTotalTime(ctx context.Context, filter bson.M) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": filter,
			},
			{
				"$group": bson.M{
					"_id": nil,
					"time": bson.M{
						"$sum": "$duration",
					},
				},
			},
		},
	)
	if err != nil {
		return 0, err
	}

	var result struct {
		Time int `bson:"time"`
	}
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return 0, err
		}
	}

	return result.Time, nil
}
