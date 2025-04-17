package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (b *TaskBusiness) GetTaskSessions(ctx context.Context, filter *entity.TaskSessionFilter) ([]entity.TaskSession, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	var taskID *string
	if filter != nil {
		taskID = filter.TaskID
	}

	var taskSessions []entity.TaskSession
	if taskID == nil {
		// Query all task sessions of the current character
		tasks, err := b.taskRepo.Find(ctx, entity.TaskPineline{
			Filter: &entity.TaskFilter{
				CharacterID: &authSession.CurrentCharacterID,
			},
		})
		if err != nil {
			return nil, err
		}

		taskIDs := lo.Map(tasks, func(task entity.Task, _ int) string {
			return task.ID
		})

		filter.TaskIDs = taskIDs
		taskSessions, err = b.taskSessionRepo.Find(ctx, entity.TaskSessionPineline{
			Filter: filter,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Query task sessions by task ID
		err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
			{
				ID:   *taskID,
				Type: entity.EntityTypeTask,
			},
		})
		if err != nil {
			return nil, err
		}

		taskSessions, err = b.taskSessionRepo.Find(ctx, entity.TaskSessionPineline{
			Filter: filter,
		})
		if err != nil {
			return nil, err
		}
	}

	return taskSessions, nil
}
