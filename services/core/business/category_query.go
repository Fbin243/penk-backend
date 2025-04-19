package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/graphql"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"
)

func (b *CategoryBusiness) Get(ctx context.Context, filter *entity.CategoryFilter, orderBy *entity.CategoryOrderBy, limit, offset *int) ([]entity.Category, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{})
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
			ID: graphql.UnassignedID,
		},
		CharacterID: authSession.CurrentCharacterID,
	})

	return cates, err
}
