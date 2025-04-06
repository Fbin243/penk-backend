package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

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

		// Create default character for the new profile
		character, err := biz.CharacterRepo.InsertOne(ctx, &entity.Character{
			BaseEntity: &base.BaseEntity{},
			ProfileID:  profile.ID,
			Name:       profile.Name,
		})
		if err != nil {
			return nil, err
		}

		profile.CurrentCharacterID = character.ID

		// Create new fish for the new user
		err = biz.CurrencyClient.CreateFish(ctx, profile.ID)
		if err != nil {
			return nil, err
		}
	}

	authSession = &rdb.AuthSession{
		ProfileID:          profile.ID,
		FirebaseUID:        profile.FirebaseUID,
		CurrentCharacterID: profile.CurrentCharacterID,
		DeviceID:           deviceID,
	}
	err = biz.Cache.SetAuthSession(ctx, profile, authSession)
	if err != nil {
		return nil, err
	}

	return authSession, nil
}
