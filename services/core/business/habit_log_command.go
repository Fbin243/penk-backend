package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/types"
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
	habitLogs, err := b.habitLogRepo.FindByHabitID(ctx, habitLogInput.HabitID, &types.TimeFilter{
		StartTime: lo.ToPtr(utils.ResetTimeToBeginningOfDay(time.Now())),
		EndTime:   lo.ToPtr(utils.ResetTimeToBeginningOfDay(time.Now())),
	})
	if err != nil {
		return nil, err
	}

	if len(habitLogs) == 0 {
		// No habit log found, create a new one
		habitLog := &entity.HabitLog{
			BaseEntity: &base.BaseEntity{},
			HabitID:    habitLogInput.HabitID,
			Timestamp:  time.Now(),
			Value:      habitLogInput.Value,
		}

		return b.habitLogRepo.InsertOne(ctx, habitLog)
	}

	// Habit log already exists, update the existing one
	habitLog := habitLogs[0]

	if habitLogInput.Value < habitLog.Value {
		reduceAmount := int64(habitLog.Value - habitLogInput.Value)

		// Reduce timetrackings duration of this habit
		timetrackings, err := b.timetrackingRepo.FindByReferenceID(ctx, habitLog.HabitID)
		if err != nil {
			return nil, err
		}

		i := 0
		removeTimeTrackingIDs := []string{}
		for i < len(timetrackings) {
			timetracking := timetrackings[i]
			if timetracking.Duration > reduceAmount {
				timetracking.Duration -= reduceAmount
				break
			} else {
				removeTimeTrackingIDs = append(removeTimeTrackingIDs, timetracking.ID)
				reduceAmount -= timetracking.Duration
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
