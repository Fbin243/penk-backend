package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/graphql"
	"tenkhours/pkg/types"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"
)

func (b *CategoryBusiness) Get(ctx context.Context, filter *entity.CategoryFilter, orderBy *entity.CategoryOrderBy, limit, offset *int) ([]entity.Category, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = &entity.CategoryFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	cates, err := b.cateRepo.Find(ctx, entity.CategoryPipeline{
		Filter:  filter,
		OrderBy: orderBy,
		Pagination: &types.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})

	// Add the default category with id = "unassigned"
	cates = append(cates, entity.Category{
		BaseEntity: &base.BaseEntity{
			ID:        graphql.UnassignedID,
			CreatedAt: utils.Now(),
			UpdatedAt: utils.Now(),
		},
		CharacterID: authSession.CurrentCharacterID,
	})

	return cates, err
}

func (b *CategoryBusiness) Count(ctx context.Context, filter *entity.CategoryFilter) (int, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return 0, err
	}

	if filter == nil {
		filter = &entity.CategoryFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	count, err := b.cateRepo.CountByFilter(ctx, filter)

	return count, err
}
