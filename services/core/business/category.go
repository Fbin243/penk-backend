package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

type CategoryBusiness struct {
	cateRepo      ICategoryRepo
	characterRepo ICharacterRepo
}

func NewCategoryBusiness(cateRepo ICategoryRepo, characterRepo ICharacterRepo) *CategoryBusiness {
	return &CategoryBusiness{cateRepo, characterRepo}
}

func (b *CategoryBusiness) GetCategories(ctx context.Context, characterID string) ([]entity.Category, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.characterRepo.ValidateCharacter(ctx, authSession.ProfileID, characterID)
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

	err := b.characterRepo.ValidateCharacter(ctx, authSession.ProfileID, cateInput.CharacterID)
	if err != nil {
		return nil, err
	}

	cate := &entity.Category{}
	if cateInput.ID != nil {
		err := b.cateRepo.ValidateCategory(ctx, cateInput.CharacterID, *cateInput.ID)
		if err != nil {
			return nil, err
		}

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
