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
	return &SnapshotsRepo{mongodb.Collection(db.SnapshotsCollection)}
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
	for cursor.Next(ctx) {
		var snapshot Snapshot
		if err := cursor.Decode(&snapshot); err != nil {
			return nil, err
		}

		snapshot.Character.ID = snapshot.Metadata.CharacterID
		snapshot.Character.UserID = snapshot.Metadata.UserID
		snapshots = append(snapshots, snapshot)
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
	for cursor.Next(ctx) {
		var snapshot Snapshot
		if err := cursor.Decode(&snapshot); err != nil {
			return nil, err
		}

		snapshot.Character.ID = snapshot.Metadata.CharacterID
		snapshot.Character.UserID = snapshot.Metadata.UserID
		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

func (r *SnapshotsRepo) GetLatestSnapshotByCharacterID(characterID primitive.ObjectID) (*Snapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	snapshot := &Snapshot{}

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
