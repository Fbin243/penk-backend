package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

// Update the user's profile
func (biz *ProfileBusiness) UpdateProfile(ctx context.Context, input entity.ProfileInput) (*entity.Profile, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
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

	updatedProfile, err := biz.ProfileRepo.FindAndUpdateByID(ctx, profile.ID, profile)
	if err != nil {
		return nil, err
	}

	return updatedProfile, nil
}

// Delete the user's profile and all related data
func (biz *ProfileBusiness) DeleteProfile(ctx context.Context) (*entity.Profile, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	var profile *entity.Profile
	characters, err := biz.CharacterRepo.GetCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	characterIDs := lo.Map(characters, func(c entity.Character, _ int) string {
		return c.ID
	})

	// Metric
	err = biz.MetricRepo.DeleteByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	// Category
	err = biz.CategoryRepo.DeleteByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	// Goal
	err = biz.GoalRepo.DeleteByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	// Time tracking
	err = biz.TimeTrackingRepo.DeleteByCharacterIDs(ctx, characterIDs)
	if err != nil {
		return nil, err
	}

	// Habit
	habits, err := biz.HabitRepo.Find(ctx, entity.HabitPipeline{
		Filter: &entity.HabitFilter{
			CharacterIDs: characterIDs,
		},
	})
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

	// Task
	tasks, err := biz.TaskRepo.Find(ctx, entity.TaskPipeline{
		Filter: &entity.TaskFilter{
			CharacterIDs: characterIDs,
		},
	})
	if err != nil {
		return nil, err
	}

	taskIDs := lo.Map(tasks, func(t entity.Task, _ int) string {
		return t.ID
	})

	err = biz.TaskSessionRepo.DeleteByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}

	// Delete all characters
	err = biz.CharacterRepo.DeleteCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	// Delete rewards
	err = biz.RewardRepo.DeleteReward(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	// Delete the profile in database
	profile, err = biz.ProfileRepo.FindAndDeleteByID(ctx, authSession.ProfileID)
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
