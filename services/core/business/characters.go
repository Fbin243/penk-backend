package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharactersBusiness struct {
	CharactersRepo      *repo.CharactersRepo
	ProfilesRepo        *repo.ProfilesRepo
	FromCreateCharacter bool
	FromUpdateCharacter bool
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
		return nil, errors.ErrorUnauthorized
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
		return nil, errors.ErrorUnauthorized
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

	// Create custom metrics for the character
	if input.CustomMetrics != nil {
		for _, customMetric := range input.CustomMetrics {
			// Insert the character into context
			ctx := context.WithValue(ctx, CharacterKey, &character)
			biz.FromCreateCharacter = true
			biz.CreateCustomMetric(ctx, character.ID, customMetric)
		}
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
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.GetCharacterByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
	}

	if input.Name != nil {
		character.Name = *input.Name
	}

	if input.Tags != nil {
		character.Tags = input.Tags
	}

	// Update custom metrics for the character
	if input.CustomMetrics != nil {
		// Insert the character into context
		ctx := context.WithValue(ctx, CharacterKey, character)
		biz.FromUpdateCharacter = true
		for _, customMetric := range input.CustomMetrics {
			if customMetric.ID != nil {
				// Update custom metric
				biz.UpdateCustomMetric(ctx, *customMetric.ID, character.ID, customMetric)
			} else {
				// Create custom metric
				biz.CreateCustomMetric(ctx, character.ID, customMetric)
			}
		}
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
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.GetCharacterByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
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
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.GetCharacterByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
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
