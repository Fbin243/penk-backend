package repo

import (
	"context"
	"time"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoalsRepo struct {
	*db.BaseRepo[Goal]
}

func NewGoalsRepo(mongo *mongo.Database) *GoalsRepo {
	return &GoalsRepo{db.NewBaseRepo[Goal](mongo.Collection(db.GoalsCollection))}
}

func (r *GoalsRepo) GetGoalsByCharacterID(characterID primitive.ObjectID) ([]Goal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	goals := []Goal{}
	cursor, err := r.Find(ctx, primitive.M{"character_id": characterID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &goals)

	return goals, err
}
