package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (b *TaskBusiness) GetTaskSessions(ctx context.Context, taskID *string) ([]entity.TaskSession, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	var taskSessions []entity.TaskSession
	if taskID == nil {
		// Query all task sessions of the current character
		tasks, err := b.taskRepo.FindByCharacterID(ctx, authSession.CurrentCharacterID)
		if err != nil {
			return nil, err
		}

		taskIDs := lo.Map(tasks, func(task entity.Task, _ int) string {
			return task.ID
		})

		taskSessions, err = b.taskSessionRepo.FindByTaskIDs(ctx, taskIDs)
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

		taskSessions, err = b.taskSessionRepo.FindByTaskID(ctx, *taskID)
		if err != nil {
			return nil, err
		}
	}

	return taskSessions, nil
}
