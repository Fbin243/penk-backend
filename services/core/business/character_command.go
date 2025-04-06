package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"
)

func (biz *CharacterBusiness) UpsertCharacter(ctx context.Context, input entity.CharacterInput) (*entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	profile, err := biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	var character *entity.Character
	if input.ID == nil {
		// Insert new character
		charactersCount, err := biz.CharacterRepo.CountCharactersByProfileID(ctx, authSession.ProfileID)
		if err != nil {
			return nil, err
		}

		if charactersCount >= int64(utils.LimitedCharacterNumber) {
			return nil, errors.ErrLimitCharacter
		}

		character = &entity.Character{
			BaseEntity: &base.BaseEntity{},
			ProfileID:  profile.ID,
		}
	} else {
		err := biz.CharacterRepo.Exist(ctx, authSession.ProfileID, *input.ID)
		if err != nil {
			return nil, err
		}

		character, err = biz.CharacterRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}
	}

	character.Name = input.Name
	if input.ID == nil {
		character, err = biz.CharacterRepo.InsertOne(ctx, character)
		if err != nil {
			return nil, err
		}

		// Update the current character of the profile with the new character
		profile.CurrentCharacterID = character.ID
		_, err = biz.ProfileRepo.UpdateByID(ctx, profile.ID, profile)
		if err != nil {
			return nil, err
		}

		err = biz.Cache.DeleteAuthSession(ctx, authSession.FirebaseUID)
		if err != nil {
			return nil, err
		}
	} else {
		character, err = biz.CharacterRepo.UpdateByID(ctx, *input.ID, character)
		if err != nil {
			return nil, err
		}
	}

	return character, nil
}

func (biz *CharacterBusiness) DeleteCharacter(ctx context.Context, id string) (*entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := biz.CharacterRepo.Exist(ctx, authSession.ProfileID, id)
	if err != nil {
		return nil, err
	}

	// Delete all metrics | habits | tasks | categories | goals of the character
	err = biz.MetricRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = biz.CategoryRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = biz.GoalRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = biz.TimeTrackingRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	character, err := biz.CharacterRepo.DeleteCharacter(ctx, id)
	if err != nil {
		return nil, err
	}

	return character, nil
}
