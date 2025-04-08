package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

func (b *TaskBusiness) UpsertTaskSession(ctx context.Context, input *entity.TaskSessionInput) (*entity.TaskSession, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
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
		taskSession, err = b.taskSessionRepo.UpdateByID(ctx, *input.ID, taskSession)
		if err != nil {
			return nil, err
		}
	}

	return taskSession, nil
}

func (b *TaskBusiness) DeleteTaskSession(ctx context.Context, id string) (*entity.TaskSession, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
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
