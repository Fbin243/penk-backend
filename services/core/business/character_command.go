package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (biz *CharacterBusiness) UpsertCharacter(ctx context.Context, input entity.CharacterInput) (*entity.Character, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
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

		if charactersCount >= utils.LimitedCharacterNumber {
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
		_, err = biz.ProfileRepo.FindAndUpdateByID(ctx, profile.ID, profile)
		if err != nil {
			return nil, err
		}

		err = biz.Cache.DeleteAuthSession(ctx, authSession.FirebaseUID)
		if err != nil {
			return nil, err
		}
	} else {
		character, err = biz.CharacterRepo.FindAndUpdateByID(ctx, *input.ID, character)
		if err != nil {
			return nil, err
		}
	}

	return character, nil
}

func (biz *CharacterBusiness) DeleteCharacter(ctx context.Context, id string) (*entity.Character, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	err = biz.CharacterRepo.Exist(ctx, authSession.ProfileID, id)
	if err != nil {
		return nil, err
	}

	// Metric
	err = biz.MetricRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Category
	err = biz.CategoryRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Goal
	err = biz.GoalRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	// TimeTracking
	err = biz.TimeTrackingRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Habit
	habits, err := biz.HabitRepo.Find(ctx, entity.HabitPipeline{
		Filter: &entity.HabitFilter{
			CharacterID: &id,
		},
	})
	if err != nil {
		return nil, err
	}

	habitIDs := lo.Map(habits, func(habit entity.Habit, _ int) string {
		return habit.ID
	})

	err = biz.HabitLogRepo.DeleteByHabitIDs(ctx, habitIDs)
	if err != nil {
		return nil, err
	}

	err = biz.HabitRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Task
	tasks, err := biz.TaskRepo.Find(ctx, entity.TaskPipeline{
		Filter: &entity.TaskFilter{
			CharacterID: &id,
		},
	})
	if err != nil {
		return nil, err
	}

	taskIDs := lo.Map(tasks, func(task entity.Task, _ int) string {
		return task.ID
	})

	err = biz.TaskSessionRepo.DeleteByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}

	err = biz.TaskRepo.DeleteByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Character
	character, err := biz.CharacterRepo.DeleteCharacter(ctx, id)
	if err != nil {
		return nil, err
	}

	return character, nil
}
