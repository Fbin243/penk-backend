package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharactersBusiness struct {
	CharactersRepo *repo.CharactersRepo
	ProfilesRepo   *repo.ProfilesRepo
}

func NewCharactersBusiness(charactersRepo *repo.CharactersRepo, profilesRepo *repo.ProfilesRepo) *CharactersBusiness {
	return &CharactersBusiness{
		CharactersRepo: charactersRepo,
		ProfilesRepo:   profilesRepo,
	}
}

func (biz *CharactersBusiness) GetCharacterByID(ctx context.Context, id string) (*repo.Character, error) {
	characterOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	character, err := biz.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to find character: %v", err)
	}

	return character, nil
}

func (biz *CharactersBusiness) GetCharactersByProfileID(ctx context.Context) ([]repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characters, err := biz.CharactersRepo.GetCharactersByProfileID(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	return characters, nil
}

func (biz *CharactersBusiness) CreateCharacter(ctx context.Context, input model.CharacterInput) (*repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	// TODO: Check if the user has already created 2 characters, maybe changed later
	characters, err := biz.CharactersRepo.GetCharactersByProfileID(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	if len(characters) >= 2 {
		return nil, fmt.Errorf("user have already created 2 characters")
	}

	character := repo.Character{
		ID:                  primitive.NewObjectID(),
		ProfileID:           profile.ID,
		TotalFocusedTime:    0,
		CustomMetrics:       []repo.CustomMetric{},
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

	createdCharacter, err := biz.CharactersRepo.CreateCharacter(&character)
	if err != nil {
		return nil, fmt.Errorf("failed to create character: %v", err)
	}

	// TODO: Character has been created, so set the current character of the user to it
	profile.CurrentCharacterID = createdCharacter.ID
	_, err = biz.ProfilesRepo.UpdateProfile(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}

	return createdCharacter, nil
}

func (biz *CharactersBusiness) UpdateCharacter(ctx context.Context, id primitive.ObjectID, input model.CharacterInput) (*repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.GetCharacterByID(id)
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

	updatedCharacter, err := biz.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return updatedCharacter, nil
}

func (biz *CharactersBusiness) DeleteCharacter(ctx context.Context, id primitive.ObjectID) (*repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.GetCharacterByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	deletedCharacter, err := biz.CharactersRepo.DeleteCharacter(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete character: %v", err)
	}

	return deletedCharacter, nil
}

func (biz *CharactersBusiness) ResetCharacter(ctx context.Context, id primitive.ObjectID) (*repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.GetCharacterByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	character.Tags = []string{}
	character.TotalFocusedTime = 0
	character.CustomMetrics = []repo.CustomMetric{}

	resetCharacter, err := biz.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to reset character: %v", err)
	}

	return resetCharacter, nil
}
