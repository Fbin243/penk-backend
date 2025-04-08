package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

func (b *TaskBusiness) GetTasks(ctx context.Context) ([]entity.Task, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	tasks, err := b.taskRepo.FindByCharacterID(ctx, authSession.CurrentCharacterID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
