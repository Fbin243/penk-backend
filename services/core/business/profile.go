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
)

type ProfileBusiness struct {
	ProfileRepo    IProfileRepo
	CharacterRepo  ICharacterRepo
	CategoryRepo   ICategoryRepo
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
		return nil, errors.ErrUnauthorized
	}

	return biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
}

// Update the user's profile
func (biz *ProfileBusiness) UpdateProfile(ctx context.Context, input entity.ProfileInput) (*entity.Profile, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	profile, err := biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	// Update the profile with the new input
	profile.Name = input.Name
	profile.ImageURL = input.ImageURL
	if input.CurrentCharacterID != nil {
		profile.CurrentCharacterID = input.CurrentCharacterID
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
		return nil, errors.ErrUnauthorized
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
	// err = auth.DeleteProfileOnFirebase(profile.FirebaseUID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to delete user profile in Firebase: %v", err)
	// }

	return profile, nil
}

func (biz *ProfileBusiness) IntrospectToken(ctx context.Context, token, deviceID string) (*rdb.AuthSession, error) {
	firebaseProfile, err := auth.GetProfileByIDToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	authSession, err := biz.Cache.GetAuthSession(ctx, firebaseProfile.UID)
	if err != nil && err != errors.ErrRedisNotFound {
		return nil, err
	}

	// Cache hit, return the profile from redis
	if authSession != nil {
		return authSession, err
	}

	// Cache misss, create a new session from the profile in DB
	profile, err := biz.ProfileRepo.GetProfileByFirebaseUID(ctx, firebaseProfile.UID)
	if err != nil && err != errors.ErrMongoNotFound {
		return nil, err
	}

	if profile == nil {
		// profile not found, mean the new account
		newProfile := entity.Profile{
			BaseEntity:  &base.BaseEntity{},
			Name:        firebaseProfile.Name,
			Email:       firebaseProfile.Email,
			FirebaseUID: firebaseProfile.UID,
			ImageURL:    firebaseProfile.Picture,
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
	}

	authSession = &rdb.AuthSession{
		ProfileID: profile.ID,
		DeviceID:  deviceID,
	}
	err = biz.Cache.SetAuthSession(ctx, profile, authSession)
	if err != nil {
		return nil, err
	}

	return authSession, nil
}

func (biz *ProfileBusiness) CheckPermission(ctx context.Context, profileID, characterID, categoryID *string) error {
	err := biz.CharacterRepo.ValidateCharacter(ctx, *profileID, *characterID)
	if err != nil {
		return err
	}

	if categoryID != nil {
		err := biz.CategoryRepo.ValidateCategory(ctx, *characterID, *categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}
