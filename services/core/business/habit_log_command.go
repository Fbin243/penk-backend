package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
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

	habitLog := &entity.HabitLog{
		BaseEntity: &base.BaseEntity{},
		HabitID:    habitLogInput.HabitID,
		Timestamp:  time.Now().UTC(),
		Value:      habitLogInput.Value,
	}

	err = b.habitLogRepo.UpsertByTimestamp(ctx, habitLog.Timestamp, habitLog)
	if err != nil {
		return nil, err
	}

	return habitLog, nil
}
