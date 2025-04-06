package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

func (biz *ProfileBusiness) GetProfile(ctx context.Context) (*entity.Profile, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	return biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
}
