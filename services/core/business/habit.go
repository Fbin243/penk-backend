package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

type HabitBusiness struct {
	habitRepo     IHabitRepo
	habitLogRepo  IHabitLogRepo
	characterRepo ICharacterRepo
	cateRepo      ICategoryRepo
}

func NewHabitBusiness(habitRepo IHabitRepo, habitLogRepo IHabitLogRepo, characterRepo ICharacterRepo, cateRepo ICategoryRepo) *HabitBusiness {
	return &HabitBusiness{
		habitRepo:     habitRepo,
		habitLogRepo:  habitLogRepo,
		characterRepo: characterRepo,
		cateRepo:      cateRepo,
	}
}

func (b *HabitBusiness) GetHabits(ctx context.Context, characterID string) ([]entity.Habit, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.characterRepo.Exist(ctx, authSession.ProfileID, characterID)
	if err != nil {
		return nil, err
	}

	return b.habitRepo.FindByCharacterID(ctx, characterID)
}

func (b *HabitBusiness) UpsertHabit(ctx context.Context, habitInput *entity.HabitInput) (*entity.Habit, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.characterRepo.Exist(ctx, authSession.ProfileID, habitInput.CharacterID)
	if err != nil {
		return nil, err
	}

	if habitInput.CategoryID != nil {
		err = b.cateRepo.Exist(ctx, habitInput.CharacterID, *habitInput.CategoryID)
		if err != nil {
			return nil, err
		}
	}

	habit := &entity.Habit{
		BaseEntity: &base.BaseEntity{},
	}
	if habitInput.ID == nil {
		count, err := b.habitRepo.CountByCharacterID(ctx, habitInput.CharacterID)
		if err != nil {
			return nil, err
		}
		if count >= utils.LimitedHabitNumber {
			return nil, errors.ErrLimitHabit
		}
	} else {
		err := b.habitRepo.Exist(ctx, habitInput.CharacterID, *habitInput.ID)
		if err != nil {
			return nil, err
		}

		habit, err = b.habitRepo.FindByID(ctx, *habitInput.ID)
		if err != nil {
			return nil, err
		}
	}

	err = copier.Copy(habit, habitInput)
	if err != nil {
		return nil, err
	}

	if habitInput.ID == nil {
		habit, err = b.habitRepo.InsertOne(ctx, habit)
		if err != nil {
			return nil, err
		}
	} else {
		habit, err = b.habitRepo.UpdateByID(ctx, *habitInput.ID, habit)
		if err != nil {
			return nil, err
		}
	}

	return habit, nil
}

func (b *HabitBusiness) DeleteHabit(ctx context.Context, habitID string) (*entity.Habit, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	habit, err := b.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	err = b.characterRepo.Exist(ctx, authSession.ProfileID, habit.CharacterID)
	if err != nil {
		return nil, err
	}

	return b.habitRepo.FindOneAndDeleteByID(ctx, habitID)
}
