package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

func (b *HabitBusiness) GetHabits(ctx context.Context) ([]entity.Habit, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	return b.habitRepo.FindByCharacterID(ctx, authSession.CurrentCharacterID)
}
