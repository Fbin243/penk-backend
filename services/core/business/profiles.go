package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"
)

type ProfilesHandler struct {
	ProfilesRepo *repo.ProfilesRepo
}

func NewProfilesHandler(profilesRepo *repo.ProfilesRepo) *ProfilesHandler {
	return &ProfilesHandler{
		ProfilesRepo: profilesRepo,
	}
}

func (r *ProfilesHandler) GetProfileByToken(ctx context.Context) (*repo.Profile, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	return &profile, nil
}

func (r *ProfilesHandler) UpdateProfile(ctx context.Context, input model.ProfileInput) (*repo.Profile, error) {
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

	updatedProfile, err := r.ProfilesRepo.UpdateProfile(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}

	return updatedProfile, nil
}
