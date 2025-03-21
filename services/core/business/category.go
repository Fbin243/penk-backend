package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"
)

type CategoryBusiness struct {
	cateRepo      ICategoryRepo
	characterRepo ICharacterRepo
	metricRepo    IMetricRepo
}

func NewCategoryBusiness(cateRepo ICategoryRepo, characterRepo ICharacterRepo, metricRepo IMetricRepo) *CategoryBusiness {
	return &CategoryBusiness{cateRepo, characterRepo, metricRepo}
}

func (b *CategoryBusiness) GetCategories(ctx context.Context, characterID string) ([]entity.Category, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.characterRepo.Exist(ctx, authSession.ProfileID, characterID)
	if err != nil {
		return nil, err
	}

	return b.cateRepo.FindByCharacterID(ctx, characterID)
}

func (b *CategoryBusiness) UpsertCategory(ctx context.Context, cateInput entity.CategoryInput) (*entity.Category, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.characterRepo.Exist(ctx, authSession.ProfileID, cateInput.CharacterID)
	if err != nil {
		return nil, err
	}

	cate := &entity.Category{}
	if cateInput.ID == nil {
		count, err := b.cateRepo.CountByCharacterID(ctx, cateInput.CharacterID)
		if err != nil {
			return nil, err
		}
		if count >= utils.LimitedCategoryNumber {
			return nil, errors.ErrLimitMetric
		}
	} else {
		cate, err = b.cateRepo.FindByID(ctx, *cateInput.ID)
		if err != nil {
			return nil, err
		}
	}

	cate.Name = cateInput.Name
	cate.CharacterID = cateInput.CharacterID
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
		return b.cateRepo.UpdateByID(ctx, *cateInput.ID, cate)
	}

	return b.cateRepo.InsertOne(ctx, cate)
}

func (b *CategoryBusiness) DeleteCategory(ctx context.Context, categoryID string) (*entity.Category, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	cate, err := b.cateRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	err = b.characterRepo.Exist(ctx, authSession.ProfileID, cate.CharacterID)
	if err != nil {
		return nil, err
	}

	// Unassign all metrics | habits | tasks of this category
	err = b.metricRepo.UnassignCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return b.cateRepo.DeleteByID(ctx, categoryID)
}
