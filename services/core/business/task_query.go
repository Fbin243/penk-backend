package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"
)

func (b *TaskBusiness) Get(ctx context.Context, filter *entity.TaskFilter, orderBy *entity.TaskOrderBy, limit, offset *int) ([]entity.Task, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	filter.CharacterID = &authSession.CurrentCharacterID
	tasks, err := b.taskRepo.Find(ctx, entity.TaskPineline{
		Filter:  filter,
		OrderBy: orderBy,
		Pagination: &types.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
