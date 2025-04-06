package business

import (
	"context"
	"math"
	"time"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/teambition/rrule-go"
)

func (b *HabitBusiness) GetHabitLogs(ctx context.Context, habitID string, startTime, endTime time.Time) ([]entity.HabitLog, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   habitID,
			Type: entity.EntityTypeHabit,
		},
	})
	if err != nil {
		return nil, err
	}

	habit, err := b.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	habitLogs, err := b.habitLogRepo.FindByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	// Parse the RRULE from string
	r, _ := rrule.StrToRRule(habit.RRule)
	if r.Options.Dtstart.Before(startTime) {
		r.DTStart(startTime)
	}

	r.Until(endTime)
	occurences := r.All()
	utils.PrintTimeSlice(occurences)

	occurenceMap := map[string]bool{}
	for _, occurence := range occurences {
		occurenceMap[occurence.Format(time.DateOnly)] = true
	}

	occurHabitLogs := []entity.HabitLog{}
	for _, habitLog := range habitLogs {
		if _, ok := occurenceMap[habitLog.Timestamp.Format(time.DateOnly)]; ok {
			if habit.Value == 0 {
				habitLog.Percent = 0
			} else {
				habitLog.Percent = math.Min(habitLog.Value/habit.Value, 1)
			}

			occurHabitLogs = append(occurHabitLogs, habitLog)
		}
	}

	return occurHabitLogs, nil
}
