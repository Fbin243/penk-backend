package business

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfilesBusiness struct {
	ProfilesRepo *repo.ProfilesRepo
	RedisClient  *redis.Client
}

func NewProfilesBusiness(profilesRepo *repo.ProfilesRepo, redisClient *redis.Client) *ProfilesBusiness {
	return &ProfilesBusiness{
		ProfilesRepo: profilesRepo,
		RedisClient:  redisClient,
	}
}

func (biz *ProfilesBusiness) GetProfile(ctx context.Context) (*repo.Profile, error) {
	// Get Firebase profile from the context
	firebaseProfile, ok := ctx.Value(auth.FirebaseProfileKey).(auth.FirebaseProfile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// Check if this UID is in Redis
	keyFound, err := biz.RedisClient.Exists(ctx, firebaseProfile.UID).Result()
	if err != nil {
		return nil, err
	}
	// Cache hit, return the profile from redis
	if keyFound == 1 {
		fmt.Print("key found")
		profileJSON, err := biz.RedisClient.Get(ctx, firebaseProfile.UID).Result()
		if err != nil {
			return nil, fmt.Errorf("login session not found in redis")
		}

		var profile repo.Profile
		err = json.Unmarshal([]byte(profileJSON), &profile)
		if err != nil {
			return nil, fmt.Errorf("failed to decode profile in redis")
		}

		return &profile, nil
	}

	fmt.Print("key not found")
	// Cache miss, create a new session from the profile in DB
	profile, err := biz.ProfilesRepo.GetProfileByFirebaseUID(firebaseProfile.UID)
	if err == mongo.ErrNoDocuments {
		// profile not found, mean the new account
		newProfile := repo.Profile{
			ID:                     primitive.NewObjectID(),
			Name:                   firebaseProfile.Name,
			Email:                  firebaseProfile.Email,
			FirebaseUID:            firebaseProfile.UID,
			ImageURL:               "",
			CreatedAt:              utils.Now(),
			UpdatedAt:              utils.Now(),
			AutoSnapshot:           true,
			AvailableSnapshots:     2,
			LimitedCharacterNumber: 2,
		}

		// Create new profile for the new user in DB
		profile, err = biz.ProfilesRepo.CreateNewProfile(&newProfile)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Save profile in redis
	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize profile: %v", err)
	}

	err = biz.RedisClient.Set(ctx, firebaseProfile.UID, profileJSON, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (biz *ProfilesBusiness) UpdateProfile(ctx context.Context, input model.ProfileInput) (*repo.Profile, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	if input.Name != nil {
		profile.Name = *input.Name
	}
	if input.ImageURL != nil {
		profile.ImageURL = *input.ImageURL
	}
	if input.CurrentCharacterID != nil {
		profile.CurrentCharacterID = *input.CurrentCharacterID
	}
	if input.AutoSnapshot != nil {
		profile.AutoSnapshot = *input.AutoSnapshot
	}

	profile.UpdatedAt = utils.Now()

	updatedProfile, err := biz.ProfilesRepo.UpdateProfile(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}

	fmt.Println("Updated profile: ", updatedProfile)
	return updatedProfile, nil
}
