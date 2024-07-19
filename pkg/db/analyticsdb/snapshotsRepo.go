package analyticsdb

import (
	"context"
	"time"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SnapshotsRepo struct {
	*mongo.Collection
}

func NewSnapshotRepo(mongodb *mongo.Database) *SnapshotsRepo {
	return &SnapshotsRepo{mongodb.Collection(db.AnalyticsCollection)}
}

func (r *SnapshotsRepo) GetSnapshotsByUserID(userID primitive.ObjectID) ([]Snapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"metadata.user_id": userID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var snapshots []Snapshot
	err = cursor.All(ctx, &snapshots)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (r *SnapshotsRepo) GetSnapshotsByCharacterID(characterID primitive.ObjectID) ([]Snapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"metadata.character_id": characterID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var snapshots []Snapshot
	err = cursor.All(ctx, &snapshots)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (r *SnapshotsRepo) GetLatestSnapshotByCharacterID(characterID primitive.ObjectID) (*Snapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var snapshot *Snapshot

	opts := options.FindOne().SetSort(bson.M{"timestamp": -1})
	err := r.FindOne(ctx, bson.M{"metadata.character_id": characterID}, opts).Decode(snapshot)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (r *SnapshotsRepo) CreateSnapshot(snapshot *Snapshot) (*Snapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, snapshot)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}
