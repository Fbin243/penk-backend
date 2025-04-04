package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/graphql"
	"tenkhours/services/core/entity"
)

func (b *CategoryBusiness) GetCategories(ctx context.Context) ([]entity.Category, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{})
	if err != nil {
		return nil, err
	}

	cates, err := b.cateRepo.FindByCharacterID(ctx, authSession.CurrentCharacterID)

	// Add the default category with id = "unassigned"
	cates = append(cates, entity.Category{
		BaseEntity: &base.BaseEntity{
			ID: graphql.UnassignedID,
		},
		CharacterID: authSession.CurrentCharacterID,
	})

	return cates, err
}
