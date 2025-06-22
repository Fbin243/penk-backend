package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (b *HabitBusiness) UpsertHabitLog(ctx context.Context, habitLogInput *entity.HabitLogInput) (*entity.HabitLog, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   habitLogInput.HabitID,
			Type: entity.EntityTypeHabit,
		},
	})
	if err != nil {
		return nil, err
	}

	if habitLogInput.Value < 0 {
		return nil, errors.NewGQLError(errors.ErrCodeBadRequest, "habit log value must be greater than or equal to 0")
	}

	habit, err := b.habitRepo.FindByID(ctx, habitLogInput.HabitID)
	if err != nil {
		return nil, err
	}

	// Check if the habit log already exists for timestamp
	timestamp, err := time.Parse(time.DateOnly, habitLogInput.Timestamp)
	if err != nil {
		return nil, errors.NewGQLError(errors.ErrCodeBadRequest, "invalid timestamp format")
	}

	habitLogs, err := b.habitLogRepo.Find(ctx, entity.HabitLogPipeline{
		Filter: &entity.HabitLogFilter{
			HabitID:   &habitLogInput.HabitID,
			StartTime: lo.ToPtr(utils.StartOfDay(timestamp)),
			EndTime:   lo.ToPtr(utils.EndOfDay(timestamp)),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(habitLogs) == 0 {
		// No habit log found, create a new one
		habitLog := &entity.HabitLog{
			BaseEntity: &base.BaseEntity{},
			HabitID:    habitLogInput.HabitID,
			Timestamp:  habitLogInput.Timestamp,
			Value:      habitLogInput.Value,
		}

		return b.habitLogRepo.InsertOne(ctx, habitLog)
	}

	// Habit log already exists, update the existing one
	habitLog := habitLogs[0]

	// Reduce timetrackings duration of this habit
	if habit.CompletionType == entity.CompletionTypeTime &&
		habitLogInput.Value < habitLog.Value {
		reduceAmount := int(habitLog.Value - habitLogInput.Value)
		timetrackings, err := b.timetrackingRepo.FindByReferenceID(ctx, habitLog.HabitID)
		if err != nil {
			return nil, err
		}

		i := 0
		removeTimeTrackingIDs := []string{}
		for i < len(timetrackings) {
			if timetrackings[i].Duration > reduceAmount {
				timetrackings[i].Duration -= reduceAmount
				break
			} else {
				removeTimeTrackingIDs = append(removeTimeTrackingIDs, timetrackings[i].ID)
				reduceAmount -= timetrackings[i].Duration
			}
			i++
		}

		// Update timetrackings
		err = b.timetrackingRepo.DeleteByIDs(ctx, removeTimeTrackingIDs)
		if err != nil {
			return nil, err
		}

		if i < len(timetrackings) {
			err = b.timetrackingRepo.UpdateByID(ctx, timetrackings[i].ID, &timetrackings[i])
			if err != nil {
				return nil, err
			}
		}
	}

	habitLog.Value = habitLogInput.Value
	return b.habitLogRepo.FindAndUpdateByID(ctx, habitLog.ID, &habitLog)
}
