package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/services/core/entity"
)

func (biz *ProfileBusiness) GetProfile(ctx context.Context) (*entity.Profile, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	return biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
}
