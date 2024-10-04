package core

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/core/validations"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfilesHandler struct {
	ProfilesRepo *coredb.ProfilesRepo
}

func NewProfilesHandler(profilesRepo *coredb.ProfilesRepo) *ProfilesHandler {
	return &ProfilesHandler{
		ProfilesRepo: profilesRepo,
	}
}

func (r *ProfilesHandler) GetProfileByToken(ctx context.Context) (*coredb.Profile, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	return &profile, nil
}

func (r *ProfilesHandler) UpdateProfile(ctx context.Context, input model.ProfileInput) (*coredb.Profile, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
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
		currentCharacterOID, err := primitive.ObjectIDFromHex(*input.CurrentCharacterID)
		if err != nil {
			return nil, err
		}

		profile.CurrentCharacterID = currentCharacterOID
	}
	if input.AutoSnapshot != nil {
		profile.AutoSnapshot = *input.AutoSnapshot
	}

	profile.UpdatedAt = utils.Now()

	err := validations.ValidateProfile(profile)
	if err != nil {
		return nil, err
	}

	updatedProfile, err := r.ProfilesRepo.UpdateProfile(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}

	return updatedProfile, nil
}
