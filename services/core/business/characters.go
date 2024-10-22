package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/coredb"
	"tenkhours/services/core/business/validations"
	"tenkhours/services/core/graph/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharactersHandler struct {
	CharactersRepo *coredb.CharactersRepo
	ProfilesRepo   *coredb.ProfilesRepo
}

func NewCharactersHandler(charactersRepo *coredb.CharactersRepo, profilesRepo *coredb.ProfilesRepo) *CharactersHandler {
	return &CharactersHandler{
		CharactersRepo: charactersRepo,
		ProfilesRepo:   profilesRepo,
	}
}

func (r *CharactersHandler) GetCharacterByID(ctx context.Context, id string) (*coredb.Character, error) {
	characterOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to find character: %v", err)
	}

	return character, nil
}

func (r *CharactersHandler) GetCharactersByProfileID(ctx context.Context) ([]coredb.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characters, err := r.CharactersRepo.GetCharactersByProfileID(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	return characters, nil
}

func (r *CharactersHandler) CreateCharacter(ctx context.Context, input model.CharacterInput) (*coredb.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	// TODO: Check if the user has already created 2 characters, maybe changed later
	characters, err := r.CharactersRepo.GetCharactersByProfileID(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	if len(characters) >= 2 {
		return nil, fmt.Errorf("user have already created 2 characters")
	}

	character := coredb.Character{
		ID:                  primitive.NewObjectID(),
		ProfileID:           profile.ID,
		TotalFocusedTime:    0,
		CustomMetrics:       []coredb.CustomMetric{},
		LimitedMetricNumber: 2,
	}

	if input.Name != nil {
		character.Name = *input.Name
	}
	if input.Gender != nil {
		character.Gender = *input.Gender
	}
	if input.Tags != nil {
		character.Tags = input.Tags
	}

	err = validations.ValidateCharacter(character)
	if err != nil {
		return nil, err
	}

	createdCharacter, err := r.CharactersRepo.CreateCharacter(&character)
	if err != nil {
		return nil, fmt.Errorf("failed to create character: %v", err)
	}

	// TODO: Character has been created, so set the current character of the user to it
	profile.CurrentCharacterID = createdCharacter.ID
	_, err = r.ProfilesRepo.UpdateProfile(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}

	return createdCharacter, nil
}

func (r *CharactersHandler) UpdateCharacter(ctx context.Context, id primitive.ObjectID, input model.CharacterInput) (*coredb.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	if input.Name != nil {
		character.Name = *input.Name
	}

	if input.Tags != nil {
		character.Tags = input.Tags
	}

	err = validations.ValidateCharacter(*character)
	if err != nil {
		return nil, err
	}

	updatedCharacter, err := r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return updatedCharacter, nil
}

func (r *CharactersHandler) DeleteCharacter(ctx context.Context, id primitive.ObjectID) (*coredb.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	deletedCharacter, err := r.CharactersRepo.DeleteCharacter(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete character: %v", err)
	}

	return deletedCharacter, nil
}

func (r *CharactersHandler) ResetCharacter(ctx context.Context, id primitive.ObjectID) (*coredb.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	character.Tags = []string{}
	character.TotalFocusedTime = 0
	character.CustomMetrics = []coredb.CustomMetric{}

	resetCharacter, err := r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to reset character: %v", err)
	}

	return resetCharacter, nil
}
