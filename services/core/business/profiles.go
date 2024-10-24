package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"
)

type ProfilesBusiness struct {
	ProfilesRepo *repo.ProfilesRepo
}

func NewProfilesBusiness(profilesRepo *repo.ProfilesRepo) *ProfilesBusiness {
	return &ProfilesBusiness{
		ProfilesRepo: profilesRepo,
	}
}

func (biz *ProfilesBusiness) GetProfileByToken(ctx context.Context) (*repo.Profile, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	return &profile, nil
}

func (biz *ProfilesBusiness) UpdateProfile(ctx context.Context, input model.ProfileInput) (*repo.Profile, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	if input.Name != nil {
		profile.Name = *input.Name
	}
	if input.ImageURL != nil {
		profile.ImageURL = *input.ImageURL
	}
	if input.CurrentCharacterID != nil {
		profile.CurrentCharacterID = *input.CurrentCharacterID
	}
	if input.AutoSnapshot != nil {
		profile.AutoSnapshot = *input.AutoSnapshot
	}

	profile.UpdatedAt = utils.Now()

	updatedProfile, err := biz.ProfilesRepo.UpdateProfile(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}

	return updatedProfile, nil
}
