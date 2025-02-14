package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileBusiness struct {
	ProfileRepo    IProfileRepo
	CharacterRepo  ICharacterRepo
	CurrencyClient ICurrencyClient
	AnalyticClient IAnalyticClient
	Cache          ICache
}

func NewProfileBusiness(profileRepo IProfileRepo, characterRepo ICharacterRepo, currencyClient ICurrencyClient, analyticClient IAnalyticClient, cache ICache) *ProfileBusiness {
	return &ProfileBusiness{
		ProfileRepo:    profileRepo,
		CharacterRepo:  characterRepo,
		CurrencyClient: currencyClient,
		AnalyticClient: analyticClient,
		Cache:          cache,
	}
}

// Get the user's profile if it exists, otherwise create a new profile
func (biz *ProfileBusiness) GetProfile(ctx context.Context) (*entity.Profile, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	return biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
}

// Update the user's profile
func (biz *ProfileBusiness) UpdateProfile(ctx context.Context, input entity.ProfileInput) (*entity.Profile, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	profile, err := biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
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

	updatedProfile, err := biz.ProfileRepo.UpdateByID(ctx, profile.ID, profile)
	if err != nil {
		return nil, err
	}

	return updatedProfile, nil
}

// Delete the user's profile and all related data
func (biz *ProfileBusiness) DeleteProfile(ctx context.Context) (*entity.Profile, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	var profile *entity.Profile
	// Delete all characters in database
	err := biz.CharacterRepo.DeleteCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete characters: %v", err)
	}

	// Delete the profile in database
	profile, err = biz.ProfileRepo.DeleteByID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user profile: %v", err)
	}

	// Delete all captured records
	err = biz.AnalyticClient.DeleteCapturedRecords(ctx, authSession.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete captured records: %v", err)
	}

	// Delete profile data from cache
	err = biz.Cache.DeleteProfileData(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to delete profile data from cache: %v", err)
	}

	// Delete user profile in Firebase
	err = auth.DeleteProfileOnFirebase(profile.FirebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user profile in Firebase: %v", err)
	}

	return profile, nil
}

func (biz *ProfileBusiness) IntrospectProfile(ctx context.Context, firebaseProfile auth.FirebaseProfile) (*entity.Profile, error) {
	profile, err := biz.ProfileRepo.GetProfileByFirebaseUID(ctx, firebaseProfile.UID)
	if err == mongo.ErrNoDocuments {
		// profile not found, mean the new account
		newProfile := entity.Profile{
			BaseEntity:         &base.BaseEntity{},
			Name:               firebaseProfile.Name,
			Email:              firebaseProfile.Email,
			FirebaseUID:        firebaseProfile.UID,
			ImageURL:           firebaseProfile.Picture,
			AutoSnapshot:       true,
			AvailableSnapshots: utils.DefaultSnapshotsNumber,
		}

		// Create new profile for the new user in DB
		profile, err = biz.ProfileRepo.InsertOne(ctx, &newProfile)
		if err != nil {
			return nil, err
		}

		// Create new fish for the new user
		err = biz.CurrencyClient.CreateFish(ctx, profile.ID)
		if err != nil {
			return nil, err
		}

	} else if err != nil {
		return nil, err
	}

	return profile, nil
}

func (biz *ProfileBusiness) CheckPermission(ctx context.Context, profileID, characterID string, metricID *string) error {
	// Check if the character belongs to the profile
	character, err := biz.CharacterRepo.FindByID(ctx, characterID)
	if err != nil {
		return err
	}

	if character.ProfileID != profileID {
		return errors.PermissionDenied()
	}

	// Check if the metric belongs to the character
	if metricID != nil {
		found := false
		for _, metric := range character.Categories {
			if metric.ID == *metricID {
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
