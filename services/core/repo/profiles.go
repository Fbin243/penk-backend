package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"tenkhours/pkg/db"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfilesRepo struct {
	*mongo.Collection
	*redis.Client
}

func NewProfilesRepo(mongodb *mongo.Database, rdb *redis.Client) *ProfilesRepo {
	profilesCollection := mongodb.Collection(db.ProfileCollection)
	_, err := profilesCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "firebase_uid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		log.Println("failed to create indexes for profiles collection")
		return nil
	}

	return &ProfilesRepo{profilesCollection, rdb}
}

func (r *ProfilesRepo) GetProfileByFirebaseUID(firebaseUID string) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var profile Profile
	err := r.FindOne(ctx, bson.M{"firebase_uid": firebaseUID}).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfilesRepo) GetProfileByEmail(email string) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var profile Profile
	err := r.FindOne(ctx, bson.M{"email": email}).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfilesRepo) CreateNewProfile(profile *Profile) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, profile)

	return profile, err
}

func (r *ProfilesRepo) UpdateProfile(profile *Profile) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.FindOneAndUpdate(ctx, bson.M{"_id": profile.ID}, bson.M{"$set": profile}, db.FindOneAndUpdateOptions).Decode(profile)

	if err == nil {
		profileJSON, err := json.Marshal(profile)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize profile: %v", err)
		}

		// Get the current TTL of the profile in Redis
		ttl, err := r.TTL(ctx, profile.FirebaseUID).Result()
		if err != nil {
			return nil, err
		}

		// Update the profile in Redis with the current TTL
		err = r.Set(ctx, profile.FirebaseUID, profileJSON, ttl).Err()
		if err != nil {
			return profile, err
		}
	}

	return profile, err
}

func (r *ProfilesRepo) DeleteProfileByFirebaseUID(firebaseUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.DeleteOne(ctx, bson.M{"firebase_uid": firebaseUID})
	return err
}
