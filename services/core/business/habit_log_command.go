package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"

	"github.com/samber/lo"
	"github.com/teambition/rrule-go"
)

func (b *HabitBusiness) UpsertHabitLog(ctx context.Context, habitLogInput *entity.HabitLogInput) (*entity.HabitLog, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
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

	// Check if today matches the habit's RRule
	rule, _ := rrule.StrToRRule(habit.RRule)
	_, found := utils.FindTimestamp(rule, time.Now())
	if !found {
		return nil, errors.NewGQLError(errors.ErrCodeBadRequest, "habit log is not valid for today")
	}

	// Check if the habit log already exists for today
	habitLogs, err := b.habitLogRepo.FindByHabitID(ctx, habitLogInput.HabitID, &entity.HabitLogFilter{
		StartTime: lo.ToPtr(utils.StartOfDay(time.Now())),
		EndTime:   lo.ToPtr(utils.EndOfDay(time.Now())),
	}, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if len(habitLogs) == 0 {
		// No habit log found, create a new one
		habitLog := &entity.HabitLog{
			BaseEntity: &base.BaseEntity{},
			HabitID:    habitLogInput.HabitID,
			Timestamp:  utils.StartOfDay(time.Now()),
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
