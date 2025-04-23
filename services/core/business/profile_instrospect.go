package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

func (biz *ProfileBusiness) IntrospectUser(ctx context.Context, token, userID, deviceID string) (*rdb.AuthSession, error) {
	authSession, err := biz.Cache.GetAuthSession(ctx, userID)
	if err != nil && err != errors.ErrRedisNotFound {
		return nil, err
	}

	// Cache hit, return the profile from redis
	if authSession != nil {
		return authSession, err
	}

	// Cache misss, create a new session from the profile in DB
	profile, err := biz.ProfileRepo.GetProfileByFirebaseUID(ctx, userID)
	if err != nil && err != errors.ErrMongoNotFound {
		return nil, err
	}

	if profile == nil {
		// profile not found, mean the new account
		// decode token to get payload
		firebaseProfile, err := auth.GetProfileByIDToken(token)
		if err != nil {
			return nil, err
		}

		newProfile := entity.Profile{
			BaseEntity: &base.BaseEntity{
				ID: mongodb.GenObjectID(),
			},
			Name:        firebaseProfile.Name,
			Email:       firebaseProfile.Email,
			FirebaseUID: firebaseProfile.UID,
			ImageURL:    firebaseProfile.Picture,
		}

		// Create default character for the new profile
		character, err := biz.CharacterRepo.InsertOne(ctx, &entity.Character{
			BaseEntity: &base.BaseEntity{},
			ProfileID:  newProfile.ID,
			Name:       newProfile.Name,
		})
		if err != nil {
			return nil, err
		}

		newProfile.CurrentCharacterID = character.ID

		// Create new profile for the new user in DB
		profile, err = biz.ProfileRepo.InsertOne(ctx, &newProfile)
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
