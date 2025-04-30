package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (b *TaskBusiness) GetTaskSessions(ctx context.Context, filter *entity.TaskSessionFilter, orderBy *entity.TaskSessionOrderBy, limit, offset *int) ([]entity.TaskSession, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	processedFilter, err := b.processFilter(ctx, authSession.CurrentCharacterID, filter)
	if err != nil {
		return nil, err
	}

	return b.taskSessionRepo.Find(ctx, entity.TaskSessionPipeline{
		Filter:  processedFilter,
		OrderBy: orderBy,
		Pagination: &types.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
}

func (b *TaskBusiness) CountTaskSession(ctx context.Context, filter *entity.TaskSessionFilter) (int, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return 0, err
	}

	processedFilter, err := b.processFilter(ctx, authSession.CurrentCharacterID, filter)
	if err != nil {
		return 0, err
	}

	return b.taskSessionRepo.CountByFilter(ctx, processedFilter)
}

func (b *TaskBusiness) processFilter(
	ctx context.Context,
	characterID string,
	filter *entity.TaskSessionFilter,
) (*entity.TaskSessionFilter, error) {
	var taskID *string
	if filter != nil {
		taskID = filter.TaskID
	} else {
		filter = &entity.TaskSessionFilter{}
	}

	if taskID == nil {
		// Get task sessions of the current character
		tasks, err := b.taskRepo.Find(ctx, entity.TaskPipeline{
			Filter: &entity.TaskFilter{
				CharacterID: &characterID,
			},
		})
		if err != nil {
			return nil, err
		}
		taskIDs := lo.Map(tasks, func(task entity.Task, _ int) string {
			return task.ID
		})
		filter.TaskIDs = taskIDs
	} else {
		// Get task sessions of the task
		err := b.permBiz.CheckOwnEntities(ctx, characterID, []PermissionEntity{
			{
				ID:   *taskID,
				Type: entity.EntityTypeTask,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return filter, nil
}
