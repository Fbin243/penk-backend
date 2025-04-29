package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

func (b *HabitBusiness) Upsert(ctx context.Context, habitInput *entity.HabitInput) (*entity.Habit, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	permEntities := []PermissionEntity{}
	if habitInput.ID != nil {
		permEntities = append(permEntities, PermissionEntity{
			ID:   *habitInput.ID,
			Type: entity.EntityTypeHabit,
		})
	}

	if habitInput.CategoryID != nil {
		permEntities = append(permEntities, PermissionEntity{
			ID:   *habitInput.CategoryID,
			Type: entity.EntityTypeCategory,
		})
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, permEntities)
	if err != nil {
		return nil, err
	}

	habit := &entity.Habit{
		BaseEntity:  &base.BaseEntity{},
		CharacterID: authSession.CurrentCharacterID,
	}
	if habitInput.ID == nil {
		count, err := b.habitRepo.CountByCharacterID(ctx, authSession.CurrentCharacterID)
		if err != nil {
			return nil, err
		}
		if count >= utils.LimitedHabitNumber {
			return nil, errors.ErrLimitHabit
		}
	} else {
		habit, err = b.habitRepo.FindByID(ctx, *habitInput.ID)
		if err != nil {
			return nil, err
		}

		if habit.CategoryID != habitInput.CategoryID {
			// Update timetrack of habit with new category
			err := b.timetrackingRepo.UpdateCategoryByReferenceID(ctx, habit.ID, habitInput.CategoryID)
			if err != nil {
				return nil, err
			}
		}

		if habit.CompletionType != habitInput.CompletionType ||
			habit.RRule != habitInput.RRule {
			// Remove all habit logs of this habit
			err := b.habitLogRepo.DeleteByHabitID(ctx, habit.ID)
			if err != nil {
				return nil, err
			}
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
		habit, err = b.habitRepo.FindAndUpdateByID(ctx, *habitInput.ID, habit)
		if err != nil {
			return nil, err
		}
	}

	return habit, nil
}

func (b *HabitBusiness) Delete(ctx context.Context, habitID string) (*entity.Habit, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   habitID,
			Type: entity.EntityTypeHabit,
		},
	})
	if err != nil {
		return nil, err
	}

	// Delete all habit logs
	err = b.habitLogRepo.DeleteByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	// Unassign reference of timetracking
	err = b.timetrackingRepo.UnassignReference(ctx, habitID, entity.EntityTypeHabit)
	if err != nil {
		return nil, err
	}

	return b.habitRepo.FindAndDeleteByID(ctx, habitID)
}
