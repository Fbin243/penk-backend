package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

func (b *TaskBusiness) GetTasks(ctx context.Context, filter *entity.TaskFilter) ([]entity.Task, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	filter.CharacterID = &authSession.CurrentCharacterID
	tasks, err := b.taskRepo.Find(ctx, entity.TaskPineline{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
