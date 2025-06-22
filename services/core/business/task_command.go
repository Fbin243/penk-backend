package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
)

func (b *TaskBusiness) Upsert(ctx context.Context, input *entity.TaskInput) (*entity.Task, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	permEntities := []PermissionEntity{}
	if input.ID != nil {
		permEntities = append(permEntities, PermissionEntity{
			ID:   *input.ID,
			Type: entity.EntityTypeTask,
		})
	}
	if input.CategoryID != nil {
		permEntities = append(permEntities, PermissionEntity{
			ID:   *input.CategoryID,
			Type: entity.EntityTypeCategory,
		})
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, permEntities)
	if err != nil {
		return nil, err
	}

	task := &entity.Task{
		BaseEntity:  &base.BaseEntity{},
		CharacterID: authSession.CurrentCharacterID,
	}

	if input.ID == nil {
		count, err := b.taskRepo.CountByCharacterID(ctx, authSession.CurrentCharacterID)
		if err != nil {
			return nil, err
		}
		if count >= utils.LimitedTaskNumber {
			return nil, errors.ErrLimitTask
		}
	} else {
		task, err = b.taskRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}

		if task.CategoryID != input.CategoryID {
			err := b.timetrackingRepo.UpdateCategoryByReferenceID(ctx, task.ID, input.CategoryID)
			if err != nil {
				return nil, err
			}
		}
	}

	err = copier.Copy(task, input)
	if err != nil {
		return nil, err
	}

	if input.Subtasks != nil {
		checkboxMap := map[string]entity.Checkbox{}
		for _, checkbox := range task.Subtasks {
			checkboxMap[checkbox.ID] = checkbox
		}

		subtasks := []entity.Checkbox{}
		for _, checkboxInput := range input.Subtasks {
			if checkboxInput.ID == nil {
				checkboxInput.ID = lo.ToPtr(mongodb.GenObjectID())
			} else {
				if _, ok := checkboxMap[*checkboxInput.ID]; !ok {
					return nil, errors.ErrBadRequest
				}

				delete(checkboxMap, *checkboxInput.ID)
			}

			checkbox := &entity.Checkbox{}
			err := copier.Copy(checkbox, checkboxInput)
			if err != nil {
				return nil, err
			}

			subtasks = append(subtasks, *checkbox)
		}

		task.Subtasks = subtasks
	}

	if input.ID == nil {
		task, err = b.taskRepo.InsertOne(ctx, task)
		if err != nil {
			return nil, err
		}
	} else {
		task, err = b.taskRepo.FindAndUpdateByID(ctx, *input.ID, task)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (b *TaskBusiness) Delete(ctx context.Context, id string) (*entity.Task, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   id,
			Type: entity.EntityTypeTask,
		},
	})
	if err != nil {
		return nil, err
	}

	// Delete all task sessions of the task
	err = b.taskSessionRepo.DeleteByTaskID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Unassign reference of timetracking
	err = b.timetrackingRepo.UnassignReference(ctx, id, entity.EntityTypeTask)
	if err != nil {
		return nil, err
	}

	task, err := b.taskRepo.FindAndDeleteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}
