package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TimeTrackingRepo) GetTotalTimeByCategoryID(ctx context.Context, categoryID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"category_id": mongodb.ToObjectID(categoryID),
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"time": bson.M{
						"$sum": bson.M{
							"$divide": []any{
								bson.M{"$subtract": []any{"$end_time", "$start_time"}},
								1000,
							},
						},
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

func (r *TimeTrackingRepo) GetTotalTimeByCharacterID(ctx context.Context, characterID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"character_id": mongodb.ToObjectID(characterID),
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"time": bson.M{
						"$sum": bson.M{
							"$divide": []any{
								bson.M{"$subtract": []any{"$end_time", "$start_time"}},
								1000,
							},
						},
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

func (r *TimeTrackingRepo) GetTotalTimeOfUnassigned(ctx context.Context, characterID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"character_id": mongodb.ToObjectID(characterID),
					"category_id":  nil,
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"time": bson.M{
						"$sum": bson.M{
							"$divide": []any{
								bson.M{"$subtract": []any{"$end_time", "$start_time"}},
								1000,
							},
						},
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
