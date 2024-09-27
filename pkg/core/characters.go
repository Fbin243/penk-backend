package core

import (
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/core/validations"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
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

func (r *CharactersHandler) GetCharacterByID(params graphql.ResolveParams) (interface{}, error) {
	profile, ok := params.Context.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	id := params.Args["id"].(string)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	return *character, nil
}

func (r *CharactersHandler) GetCharactersByProfileID(params graphql.ResolveParams) (interface{}, error) {
	profile, ok := params.Context.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characters, err := r.CharactersRepo.GetCharactersByProfileID(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	return characters, nil
}

func (r *CharactersHandler) GetAllCharacters(params graphql.ResolveParams) (interface{}, error) {
	_, ok := params.Context.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	// TODO: Check if the user is admin ...

	characters, err := r.CharactersRepo.GetAllCharacters()
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	return characters, nil
}

func (r *CharactersHandler) CreateCharacter(params graphql.ResolveParams) (interface{}, error) {
	profile, ok := params.Context.Value(auth.ProfileKey).(coredb.Profile)
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

	input := params.Args["input"].(map[string]interface{})
	if name, ok := input["name"].(string); ok {
		character.Name = name
	}

	if gender, ok := input["gender"].(bool); ok {
		character.Gender = gender
	}

	if tags, ok := input["tags"].([]interface{}); ok {
		character.Tags = convertListToSlice(tags)
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

	return *createdCharacter, nil
}

func (r *CharactersHandler) UpdateCharacter(params graphql.ResolveParams) (interface{}, error) {
	profile, ok := params.Context.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	id := params.Args["id"].(string)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	input := params.Args["input"].(map[string]interface{})
	if name, ok := input["name"].(string); ok {
		character.Name = name
	}

	// TODO: It may be added later or not :)
	// if gender, ok := input["gender"].(bool); ok {
	// 	character.Gender = gender
	// }

	if tags, ok := input["tags"].([]interface{}); ok {
		character.Tags = convertListToSlice(tags)
	}

	err = validations.ValidateCharacter(*character)
	if err != nil {
		return nil, err
	}

	updatedCharacter, err := r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return *updatedCharacter, nil
}

func (r *CharactersHandler) DeleteCharacter(params graphql.ResolveParams) (interface{}, error) {
	profile, ok := params.Context.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	id := params.Args["id"].(string)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	deletedCharacter, err := r.CharactersRepo.DeleteCharacter(objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete character: %v", err)
	}

	return *deletedCharacter, nil
}

func (r *CharactersHandler) ResetCharacter(params graphql.ResolveParams) (interface{}, error) {
	profile, ok := params.Context.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	id := params.Args["id"].(string)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if objectID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	character, err := r.CharactersRepo.GetCharacterByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	character.Tags = []string{}
	character.TotalFocusedTime = 0
	character.CustomMetrics = []coredb.CustomMetric{}

	updatedCharacter, err := r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to reset character: %v", err)
	}

	return *updatedCharacter, nil
}
