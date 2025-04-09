package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	rrulex "tenkhours/pkg/rrule"
	"tenkhours/services/core/entity"

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
	_, found := rrulex.FindTimestamp(rule, time.Now())
	if !found {
		return nil, errors.NewGQLError(errors.ErrCodeBadRequest, "habit log is not valid for today")
	}

	habitLog := &entity.HabitLog{
		BaseEntity: &base.BaseEntity{},
		HabitID:    habitLogInput.HabitID,
		Timestamp:  time.Now(),
		Value:      habitLogInput.Value,
	}

	err = b.habitLogRepo.UpsertByTimestamp(ctx, habitLog.Timestamp, habitLog)
	if err != nil {
		return nil, err
	}

	return habitLog, nil
}
