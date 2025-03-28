package business

import (
	"context"
	"math"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (b *HabitBusiness) GetHabitLogs(ctx context.Context, habitID string) ([]entity.HabitLog, error) {
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

	habitLogs, err := b.habitLogRepo.FindByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	habitLogs = lo.Map(habitLogs, func(habitLog entity.HabitLog, _ int) entity.HabitLog {
		if habit.Value == 0 {
			habitLog.Percent = 0
		} else {
			habitLog.Percent = math.Min(habitLog.Value/habit.Value, 1)
		}
		return habitLog
	})

	return habitLogs, nil
}

func (b *HabitBusiness) UpsertHabitLog(ctx context.Context, habitLogInput *entity.HabitLogInput) (*entity.HabitLog, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	habit, err := b.habitRepo.FindByID(ctx, habitLogInput.HabitID)
	if err != nil {
		return nil, err
	}

	err = b.characterRepo.Exist(ctx, authSession.ProfileID, habit.CharacterID)
	if err != nil {
		return nil, err
	}

	habitLog := &entity.HabitLog{
		BaseEntity: &base.BaseEntity{},
		HabitID:    habitLogInput.HabitID,
		Timestamp:  time.Now(),
		Value:      habitLogInput.Value,
	}

	err = b.habitLogRepo.DeleteByHabitIDAndTimestamp(ctx, habitLog.HabitID, habitLog.Timestamp)
	if err != nil {
		return nil, err
	}

	habitLog, err = b.habitLogRepo.InsertOne(ctx, habitLog)
	if err != nil {
		return nil, err
	}

	return habitLog, nil
}
