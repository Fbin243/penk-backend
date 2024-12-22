package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/utils"
	"tenkhours/services/currency/repo"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfilesRepo struct {
	*db.BaseRepo[Profile]
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

	return &ProfilesRepo{db.NewBaseRepo[Profile](profilesCollection), rdb}
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

func (r *ProfilesRepo) UpdateProfile(profile *Profile) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile, err := r.UpdateByID(profile.ID, profile)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	return profile, err
}

func (r *ProfilesRepo) DeleteProfileByFirebaseUID(firebaseUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.DeleteOne(ctx, bson.M{"firebase_uid": firebaseUID})
	return err
}

// Find profile in Redis and DB, create a new profile if not found
func (r *ProfilesRepo) FindOrCreateProfile(ctx context.Context, firebaseProfile *auth.FirebaseProfile) (*Profile, error) {
	// Check if there is any active session in Redis
	keyFound, err := r.Exists(context.Background(), firebaseProfile.UID).Result()
	if err != nil {
		return nil, err
	}

	// Cache hit, return the profile from redis
	if keyFound == 1 {
		profileJSON, err := r.Get(context.Background(), firebaseProfile.UID).Result()
		if err != nil {
			return nil, err
		}

		var profile Profile
		err = json.Unmarshal([]byte(profileJSON), &profile)
		if err != nil {
			return nil, err
		}
	}

	// Cache miss, create a new session from the profile in DB
	profile, err := r.GetProfileByFirebaseUID(firebaseProfile.UID)
	if err == mongo.ErrNoDocuments {
		// profile not found, mean the new account
		newProfile := Profile{
			BaseModel:              &db.BaseModel{},
			Name:                   firebaseProfile.Name,
			Email:                  firebaseProfile.Email,
			FirebaseUID:            firebaseProfile.UID,
			ImageURL:               firebaseProfile.Picture,
			AutoSnapshot:           true,
			AvailableSnapshots:     utils.DefaultSnapshotsNumber,
			LimitedCharacterNumber: utils.LimitedCharacterNumber,
		}

		// Create new profile for the new user in DB
		profile, err = r.InsertOne(&newProfile)
		if err != nil {
			return nil, err
		}

		// TODO: Create new fish by making temp fish repo in the profile repo @Fbin243
		fishRepo := repo.NewFishRepo(r.Database())
		newFish := &repo.Fish{
			BaseModel: &db.BaseModel{},
			ProfileID: profile.ID,
			Gold:      0,
			Normal:    0,
		}

		_, err := fishRepo.InsertOne(newFish)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Save profile in redis
	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return nil, err
	}

	err = r.Set(context.Background(), firebaseProfile.UID, profileJSON, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return profile, nil
}
