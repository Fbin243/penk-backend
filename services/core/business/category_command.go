package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"
)

func (b *CategoryBusiness) UpsertCategory(ctx context.Context, cateInput entity.CategoryInput) (*entity.Category, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	cate := &entity.Category{
		BaseEntity:  &base.BaseEntity{},
		CharacterID: authSession.CurrentCharacterID,
	}
	if cateInput.ID == nil {
		count, err := b.cateRepo.CountByCharacterID(ctx, authSession.CurrentCharacterID)
		if err != nil {
			return nil, err
		}
		if count >= utils.LimitedCategoryNumber {
			return nil, errors.ErrLimitMetric
		}
	} else {
		err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
			{
				ID:   *cateInput.ID,
				Type: entity.EntityTypeCategory,
			},
		})
		if err != nil {
			return nil, err
		}

		cate, err = b.cateRepo.FindByID(ctx, *cateInput.ID)
		if err != nil {
			return nil, err
		}
	}

	cate.Name = cateInput.Name
	if cateInput.Description != nil {
		cate.Description = *cateInput.Description
	}
	if cateInput.Style != nil {
		cate.Style = entity.CategoryStyle{
			Color: cateInput.Style.Color,
			Icon:  cateInput.Style.Icon,
		}
	}

	if cateInput.ID != nil {
		return b.cateRepo.FindAndUpdateByID(ctx, *cateInput.ID, cate)
	}

	return b.cateRepo.InsertOne(ctx, cate)
}

func (b *CategoryBusiness) DeleteCategory(ctx context.Context, categoryID string) (*entity.Category, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{})
	if err != nil {
		return nil, err
	}

	// Unassign category
	// Metric
	err = b.metricRepo.UnassignCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Timetracking
	err = b.timetrackingRepo.UnassignCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Habit
	err = b.habitRepo.UnassignCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Task
	err = b.taskRepo.UnassignCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return b.cateRepo.FindAndDeleteByID(ctx, categoryID)
}
