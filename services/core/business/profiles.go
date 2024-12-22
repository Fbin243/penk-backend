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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	return biz.ProfilesRepo.FindByID(authSession.ProfileID)
}

// Update the user's profile
func (biz *ProfilesBusiness) UpdateProfile(ctx context.Context, input model.ProfileInput) (*repo.Profile, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	profile, err := biz.ProfilesRepo.FindByID(authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	// Update the profile with the new input
	profile.Name = input.Name
	profile.ImageURL = input.ImageURL
	if input.CurrentCharacterID != nil {
		profile.CurrentCharacterID = *input.CurrentCharacterID
	}
	if input.AutoSnapshot != nil {
		profile.AutoSnapshot = *input.AutoSnapshot
	}

	profile.UpdatedAt = utils.Now()

	updatedProfile, err := biz.ProfilesRepo.UpdateProfile(profile)
	if err != nil {
		return nil, err
	}

	return updatedProfile, nil
}

// Delete the user's profile and all related data
func (biz *ProfilesBusiness) DeleteProfile(ctx context.Context) (*repo.Profile, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	var profile *repo.Profile
	// Make a callback for deleting profile and related data
	callback := func(ctx mongo.SessionContext) (interface{}, error) {
		// Delete all captured records in database
		err := biz.CapturedRecordsRepo.DeleteCapturedRecordsByProfileID(authSession.ProfileID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete captured records: %v", err)
		}

		// Delete all snapshots in database
		err = biz.SnapshotsRepo.DeleteSnapshotsByProfileID(authSession.ProfileID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete snapshots: %v", err)
		}

		// Delete all characters in database
		err = biz.CharactersRepo.DeleteCharactersByProfileID(authSession.ProfileID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete characters: %v", err)
		}

		// Delete the profile in database
		profile, err = biz.ProfilesRepo.DeleteByID(authSession.ProfileID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete user profile: %v", err)
		}

		return profile, nil
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

	return profile, nil
}

func (biz *ProfilesBusiness) IntrospectProfile(ctx context.Context, firebaseProfile auth.FirebaseProfile) (*repo.Profile, error) {
	profile, err := biz.ProfilesRepo.GetProfileByFirebaseUID(firebaseProfile.UID)
	if err == mongo.ErrNoDocuments {
		// profile not found, mean the new account
		newProfile := repo.Profile{
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
		profile, err = biz.ProfilesRepo.InsertOne(&newProfile)
		if err != nil {
			return nil, err
		}

		// TODO: Create new fish by making temp fish repo in the profile repo @Fbin243
		_fishRepo := fishRepo.NewFishRepo(db.GetDBManager().DB)
		newFish := &fishRepo.Fish{
			BaseModel: &db.BaseModel{},
			ProfileID: profile.ID,
			Gold:      0,
			Normal:    0,
		}

		_, err := _fishRepo.InsertOne(newFish)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return profile, nil
}

func (biz *ProfilesBusiness) CheckPermission(ctx context.Context, profileID, characterID, metricID primitive.ObjectID) error {
	// Check if the character belongs to the profile
	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return err
	}

	if character.ProfileID != profileID {
		return errors.PermissionDenied()
	}

	// Check if the metric belongs to the character
	if !metricID.IsZero() {
		found := false
		for _, metric := range character.CustomMetrics {
			if metric.ID == metricID {
				found = true
				break
			}
		}

		if !found {
			return errors.PermissionDenied()
		}
	}

	return nil
}
