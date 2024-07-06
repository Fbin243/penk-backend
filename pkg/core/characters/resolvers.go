package characters

import (
	"fmt"
	"log"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharactersResolver struct {
	CharactersRepo *coredb.CharactersRepo
	UsersRepo      *coredb.UsersRepo
}

func NewCharactersResolver(charactersRepo *coredb.CharactersRepo, usersRepo *coredb.UsersRepo) *CharactersResolver {
	return &CharactersResolver{
		CharactersRepo: charactersRepo,
		UsersRepo:      usersRepo,
	}
}

func (r *CharactersResolver) GetCharacterByID(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
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

	if character.UserID != user.ID {
		return nil, fmt.Errorf("permission denied")
	}

	return *character, nil
}

func (r *CharactersResolver) GetCharactersByUserID(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	characters, err := r.CharactersRepo.GetCharactersByUserID(user.ID)
	if err != nil {
		log.Printf("failed to find character: %v\n", err)
		return nil, err
	}

	return characters, nil
}

func (r *CharactersResolver) GetAllCharacters(params graphql.ResolveParams) (interface{}, error) {
	_, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	// TODO: Check if the user is admin ...

	characters, err := r.CharactersRepo.GetAllCharacters()
	if err != nil {
		log.Printf("failed to find characters: %v\n", err)
		return nil, err
	}

	return characters, nil
}

func (r *CharactersResolver) CreateCharacter(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	// TODO: Check if the user has already created 2 characters, maybe changed later
	characters, err := r.CharactersRepo.GetCharactersByUserID(user.ID)
	if err != nil {
		log.Printf("failed to find user's character: %v\n", err)
		return nil, err
	}

	if len(characters) >= 2 {
		return nil, fmt.Errorf("user have already created 2 characters")
	}

	character := coredb.Character{
		ID:                  primitive.NewObjectID(),
		UserID:              user.ID,
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

	if avatar, ok := input["avatar"].(string); ok {
		character.Avatar = avatar
	}

	if tags, ok := input["tags"].([]interface{}); ok {
		character.Tags = convertListToSlice(tags)
	}

	err = ValidateCharacter(character)
	if err != nil {
		return nil, err
	}

	createdCharacter, err := r.CharactersRepo.CreateCharacter(&character)
	if err != nil {
		log.Printf("failed to create character: %v\n", err)
		return nil, err
	}

	// TODO: Character has been created, so set the current character of the user to it
	user.CurrentCharacterID = createdCharacter.ID
	_, err = r.UsersRepo.UpdateUser(&user)
	if err != nil {
		log.Printf("failed to update current character: %v\n", err)
	}

	return *createdCharacter, nil
}

func (r *CharactersResolver) UpdateCharacter(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
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

	if character.UserID != user.ID {
		return nil, fmt.Errorf("permission denied")
	}

	input := params.Args["input"].(map[string]interface{})
	if name, ok := input["name"].(string); ok {
		character.Name = name
	}

	// TODO: It may be added later or not :)
	// if gender, ok := input["gender"].(bool); ok {
	// 	character.Gender = gender
	// }

	if avatar, ok := input["avatar"].(string); ok {
		character.Avatar = avatar
	}

	if tags, ok := input["tags"].([]interface{}); ok {
		character.Tags = convertListToSlice(tags)
	}

	err = ValidateCharacter(*character)
	if err != nil {
		return nil, err
	}

	updatedCharacter, err := r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return *updatedCharacter, nil
}

func (r *CharactersResolver) DeleteCharacter(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	id := params.Args["id"].(string)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if objectID != user.ID {
		return nil, fmt.Errorf("permission denied")
	}

	deletedCharacter, err := r.CharactersRepo.DeleteCharacter(objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete character: %v", err)
	}

	return *deletedCharacter, nil
}

func (r *CharactersResolver) ResetCharacter(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	id := params.Args["id"].(string)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if objectID != user.ID {
		return nil, fmt.Errorf("permission denied")
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
