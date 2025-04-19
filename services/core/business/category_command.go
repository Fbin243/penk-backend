package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"
)

func (b *CategoryBusiness) Upsert(ctx context.Context, input *entity.CategoryInput) (*entity.Category, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	cate := &entity.Category{
		BaseEntity:  &base.BaseEntity{},
		CharacterID: authSession.CurrentCharacterID,
	}
	if input.ID == nil {
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
				ID:   *input.ID,
				Type: entity.EntityTypeCategory,
			},
		})
		if err != nil {
			return nil, err
		}

		cate, err = b.cateRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}
	}

	cate.Name = input.Name
	if input.Description != nil {
		cate.Description = *input.Description
	}
	if input.Style != nil {
		cate.Style = entity.CategoryStyle{
			Color: input.Style.Color,
			Icon:  input.Style.Icon,
		}
	}

	if input.ID != nil {
		return b.cateRepo.FindAndUpdateByID(ctx, *input.ID, cate)
	}

	return b.cateRepo.InsertOne(ctx, cate)
}

func (b *CategoryBusiness) Delete(ctx context.Context, categoryID string) (*entity.Category, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{})
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
