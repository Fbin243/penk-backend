package core

import (
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/core/validations"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/graphql-go/graphql"
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

func (r *ProfilesHandler) GetProfileByToken(params graphql.ResolveParams) (interface{}, error) {
	profile, ok := params.Context.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	return profile, nil
}

func (r *ProfilesHandler) UpdateAccount(params graphql.ResolveParams) (interface{}, error) {
	profile, ok := params.Context.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	input := params.Args["input"].(map[string]interface{})
	if name, ok := input["name"].(string); ok {
		profile.Name = name
	}

	if imageURL, ok := input["imageURL"].(string); ok {
		profile.ImageURL = imageURL
	}

	if currentCharacterID, ok := input["currentCharacterID"].(string); ok {
		currentCharacterOID, err := primitive.ObjectIDFromHex(currentCharacterID)
		if err != nil {
			return nil, err
		}

		profile.CurrentCharacterID = currentCharacterOID
	}

	if autoSnapshot, ok := input["autoSnapshot"].(bool); ok {
		profile.AutoSnapshot = autoSnapshot
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

	return *updatedProfile, nil
}
