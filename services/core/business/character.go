package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"

	"github.com/samber/lo"
)

type CharacterBusiness struct {
	CharacterRepo ICharacterRepo
	ProfileRepo   IProfileRepo
	GoalRepo      IGoalRepo
}

func NewCharacterBusiness(characterRepo ICharacterRepo, profileRepo IProfileRepo, goalRepo IGoalRepo) *CharacterBusiness {
	return &CharacterBusiness{
		CharacterRepo: characterRepo,
		ProfileRepo:   profileRepo,
		GoalRepo:      goalRepo,
	}
}

func (biz *CharacterBusiness) GetCharacterByID(ctx context.Context, id string) (*entity.Character, error) {
	character, err := biz.CharacterRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (biz *CharacterBusiness) GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	characters, err := biz.CharacterRepo.GetCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

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
			return nil, errors.NewGQLError(errors.ErrCodeLimitCharacter, nil)
		}

		character = &entity.Character{
			BaseEntity: &base.BaseEntity{},
			ProfileID:  profile.ID,
		}
	} else {
		// Update existing character
		character, err = biz.CharacterRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}

		if ok, _ := auth.CheckPermission(profile, character, "write"); !ok {
			return nil, errors.ErrPermissionDenied
		}
	}

	character.Name = input.Name
	if input.ID == nil {
		character, err = biz.CharacterRepo.InsertOne(ctx, character)
		if err != nil {
			return nil, err
		}

		// Update the current character of the profile with the new character
		profile.CurrentCharacterID = lo.ToPtr(character.ID)
		_, err = biz.ProfileRepo.UpdateByID(ctx, profile.ID, profile)
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

	character, err := biz.CharacterRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	profile := &entity.Profile{
		BaseEntity: &base.BaseEntity{
			ID: authSession.ProfileID,
		},
	}

	if ok, _ := auth.CheckPermission(profile, character, "write"); !ok {
		return nil, errors.ErrPermissionDenied
	}

	deletedCharacter, err := biz.CharacterRepo.DeleteCharacter(ctx, id)
	if err != nil {
		return nil, err
	}

	return deletedCharacter, nil
}
