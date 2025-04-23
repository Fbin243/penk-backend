package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

func (b *TaskBusiness) UpsertTaskSession(ctx context.Context, input *entity.TaskSessionInput) (*entity.TaskSession, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   input.TaskID,
			Type: entity.EntityTypeTask,
		},
	})
	if err != nil {
		return nil, err
	}

	taskSession := &entity.TaskSession{
		BaseEntity: &base.BaseEntity{},
	}

	if input.ID != nil {
		taskSession, err = b.taskSessionRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}
	}

	err = copier.Copy(taskSession, input)
	if err != nil {
		return nil, err
	}

	if input.ID == nil {
		taskSession, err = b.taskSessionRepo.InsertOne(ctx, taskSession)
		if err != nil {
			return nil, err
		}
	} else {
		taskSession, err = b.taskSessionRepo.FindAndUpdateByID(ctx, *input.ID, taskSession)
		if err != nil {
			return nil, err
		}
	}

	return taskSession, nil
}

func (b *TaskBusiness) UpsertTaskSessions(ctx context.Context, inputs []entity.TaskSessionInput) ([]entity.TaskSession, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	// Check task ids of all task sessions
	updateTaskSessionIDs := []string{}
	updateTaskSessionMap := map[string]entity.TaskSessionInput{}
	insertTaskSessions := []entity.TaskSession{}
	permEntities := []PermissionEntity{}

	for _, input := range inputs {
		permEntities = append(permEntities, PermissionEntity{
			ID:   input.TaskID,
			Type: entity.EntityTypeTask,
		})
		if input.ID != nil {
			updateTaskSessionIDs = append(updateTaskSessionIDs, *input.ID)
			updateTaskSessionMap[*input.ID] = input
		} else {
			taskSession := &entity.TaskSession{
				BaseEntity: &base.BaseEntity{},
			}

			err = copier.Copy(taskSession, input)
			if err != nil {
				return nil, err
			}

			insertTaskSessions = append(insertTaskSessions, *taskSession)
		}
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, permEntities)
	if err != nil {
		return nil, err
	}

	// Insert new task sessions
	insertTaskSessions, err = b.taskSessionRepo.InsertMany(ctx, insertTaskSessions)
	if err != nil {
		return nil, err
	}

	// Update existing task sessions
	// --- Find all task sessions
	existingTaskSessions, err := b.taskSessionRepo.FindByIDs(ctx, updateTaskSessionIDs)
	if err != nil {
		return nil, err
	}

	// --- Update task sessions
	updateTaskSessions := make([]entity.TaskSession, len(existingTaskSessions))
	for _, taskSession := range existingTaskSessions {
		input := updateTaskSessionMap[taskSession.ID]
		err = copier.Copy(taskSession, input)
		if err != nil {
			return nil, err
		}

		updateTaskSessions = append(updateTaskSessions, taskSession)
	}

	// --- Update task sessions in DB
	updateTaskSessions, err = b.taskSessionRepo.FindAndUpdateByIDs(ctx, updateTaskSessions)
	if err != nil {
		return nil, err
	}

	return append(insertTaskSessions, updateTaskSessions...), nil
}

func (b *TaskBusiness) DeleteTaskSession(ctx context.Context, id string) (*entity.TaskSession, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	taskSession, err := b.taskSessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   taskSession.TaskID,
			Type: entity.EntityTypeTask,
		},
	})
	if err != nil {
		return nil, err
	}

	err = b.taskSessionRepo.DeleteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return taskSession, nil
}
