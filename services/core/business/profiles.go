package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	analyticsRepo "tenkhours/services/analytics/repo"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"
	fishRepo "tenkhours/services/currency/repo"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfilesBusiness struct {
	ProfilesRepo        *repo.ProfilesRepo
	FishRepo            *fishRepo.FishRepo
	CharactersRepo      *repo.CharactersRepo
	CapturedRecordsRepo *analyticsRepo.CapturedRecordsRepo
	SnapshotsRepo       *analyticsRepo.SnapshotsRepo
	RedisClient         *redis.Client
}

func NewProfilesBusiness(profilesRepo *repo.ProfilesRepo, fishRepo *fishRepo.FishRepo, charactersRepo *repo.CharactersRepo, capturedRecordsRepo *analyticsRepo.CapturedRecordsRepo, snapshotsRepo *analyticsRepo.SnapshotsRepo, redisClient *redis.Client) *ProfilesBusiness {
	return &ProfilesBusiness{
		ProfilesRepo:        profilesRepo,
		CharactersRepo:      charactersRepo,
		CapturedRecordsRepo: capturedRecordsRepo,
		SnapshotsRepo:       snapshotsRepo,
		FishRepo:            fishRepo,
		RedisClient:         redisClient,
	}
}

// Get the user's profile if it exists, otherwise create a new profile
func (biz *ProfilesBusiness) GetProfile(ctx context.Context) (*repo.Profile, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// TODO-NAM: Help Fbin to move this logic to the middleware, after create new profile
	// // Create new fish for user
	// newFish := fishRepo.Fish{
	// 	ID:        primitive.NewObjectID(),
	// 	ProfileID: newProfile.ID,
	// 	Gold:      0,
	// 	Normal:    0,
	// }

	// createdFish, err := biz.FishRepo.CreateFish(&newFish)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create fish for new profile: %v", err)
	// }
	// fmt.Printf("Created fish: %+v\n", createdFish)
	return &profile, nil
}

// Update the user's profile
func (biz *ProfilesBusiness) UpdateProfile(ctx context.Context, input model.ProfileInput) (*repo.Profile, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// Update the profile with the new input
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
		fmt.Println("input.AutoSnapshot", *input.AutoSnapshot)
		profile.AutoSnapshot = *input.AutoSnapshot
	}

	profile.UpdatedAt = utils.Now()

	updatedProfile, err := biz.ProfilesRepo.UpdateProfile(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}

	return updatedProfile, nil
}

// Delete the user's profile and all related data
func (biz *ProfilesBusiness) DeleteProfile(ctx context.Context) (*repo.Profile, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// Make a callback for deleting profile and related data
	callback := func(ctx mongo.SessionContext) (interface{}, error) {
		// Delete all captured records in database
		err := biz.CapturedRecordsRepo.DeleteCapturedRecordsByProfileID(profile.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete captured records: %v", err)
		}

		// Delete all snapshots in database
		err = biz.SnapshotsRepo.DeleteSnapshotsByProfileID(profile.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete snapshots: %v", err)
		}

		// Delete all characters in database
		err = biz.CharactersRepo.DeleteCharactersByProfileID(profile.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete characters: %v", err)
		}

		// Delete the profile in database
		err = biz.ProfilesRepo.DeleteProfileByFirebaseUID(profile.FirebaseUID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete user profile: %v", err)
		}

		return nil, nil
	}

	// Execute the callback in a transaction
	session, err := db.GetDBManager().Client.StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), callback)
	if err != nil {
		return nil, err
	}

	// Delete current captured record in redis
	err = biz.RedisClient.Del(ctx, db.GetCapturedRecordKey(profile.ID.Hex())).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to delete current captured record in redis: %v", err)
	}

	// Delete current timetracking in redis
	err = biz.RedisClient.Del(ctx, profile.ID.Hex()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to delete current timetracking in redis: %v", err)
	}

	// Delete profile in redis
	err = biz.RedisClient.Del(ctx, profile.FirebaseUID).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to delete profile in redis: %v", err)
	}

	// Delete user profile in Firebase
	err = auth.DeleteProfileOnFirebase(profile.FirebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user profile in Firebase: %v", err)
	}

	return &profile, nil
}
