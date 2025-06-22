package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"
)

func (b *TaskBusiness) Get(ctx context.Context, filter *entity.TaskFilter, orderBy *entity.TaskOrderBy, limit, offset *int) ([]entity.Task, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = &entity.TaskFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	return b.taskRepo.Find(ctx, entity.TaskPipeline{
		Filter:  filter,
		OrderBy: orderBy,
		Pagination: &types.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
}

func (b *TaskBusiness) Count(ctx context.Context, filter *entity.TaskFilter) (int, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return 0, err
	}

	if filter == nil {
		filter = &entity.TaskFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	return b.taskRepo.CountByFilter(ctx, filter)
}
