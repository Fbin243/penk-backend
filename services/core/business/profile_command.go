package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

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
	if profile.CurrentCharacterID != input.CurrentCharacterID {
		err := biz.permBiz.CheckOwnCharacter(ctx, profile.ID, input.CurrentCharacterID)
		if err != nil {
			return nil, err
		}

		err = biz.Cache.DeleteAuthSession(ctx, profile.FirebaseUID)
		if err != nil {
			return nil, err
		}
	}

	profile.Name = input.Name
	profile.ImageURL = input.ImageURL
	profile.CurrentCharacterID = input.CurrentCharacterID
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
	// Find all characters of profile
	characters, err := biz.CharacterRepo.GetCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	// Delete all metrics | habits | tasks | categories | goals | timetrackings of all characters
	characterIDs := lo.Map(characters, func(c entity.Character, _ int) string {
		return c.ID
	})

	err = biz.MetricRepo.DeleteByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	err = biz.CategoryRepo.DeleteByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	err = biz.GoalRepo.DeleteByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	habits, err := biz.HabitRepo.FindByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	habitIDs := lo.Map(habits, func(h entity.Habit, _ int) string {
		return h.ID
	})

	err = biz.HabitLogRepo.DeleteByHabitIDs(ctx, habitIDs)
	if err != nil {
		return nil, err
	}

	err = biz.HabitRepo.DeleteByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	err = biz.TimeTrackingRepo.DeleteByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	// Delete all characters in database
	err = biz.CharacterRepo.DeleteCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	// TODO: @Namiscrea7or refactor later after the usecase of currency is clear
	// err = biz.CurrencyClient.DeleteFish(ctx, authSession.ProfileID)
	// if err != nil {
	// 	return nil, err
	// }

	// Delete the profile in database
	profile, err = biz.ProfileRepo.FindOneAndDeleteByID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	// Delete profile data from cache
	err = biz.Cache.DeleteProfileData(ctx, profile)
	if err != nil {
		return nil, err
	}

	// Delete user profile in Firebase
	err = auth.DeleteProfileOnFirebase(profile.FirebaseUID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
