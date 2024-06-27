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
}

func NewCharactersResolver() *CharactersResolver {
	return &CharactersResolver{
		CharactersRepo: coredb.NewCharactersRepo(),
	}
}

func (r *CharactersResolver) CreateCharacter(params graphql.ResolveParams) (interface{}, error) {
	user := params.Context.Value(auth.UserKey).(*coredb.User)
	name := params.Args["name"].(string)

	var tags []string
	if tagsList, ok := params.Args["tags"].([]interface{}); ok {
		tags = convertListToSlice(tagsList)
	}

	character := coredb.Character{
		ID:                  primitive.NewObjectID(),
		UserID:              user.ID,
		Name:                name,
		Tags:                tags,
		TotalFocusedTime:    0,
		CustomMetrics:       []coredb.CustomMetric{},
		LimitedMetricNumber: 2,
	}

	createResult, err := r.CharactersRepo.CreateCharacter(character)
	if err != nil {
		log.Printf("failed to create character: %v\n", err)
		return nil, err
	}

	return createResult, nil
}

func (r *CharactersResolver) GetCharacterByID(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	return character, nil
}

func (r *CharactersResolver) GetCharactersByUserID(params graphql.ResolveParams) (interface{}, error) {
	user := params.Context.Value(auth.UserKey).(*coredb.User)

	characters, err := r.CharactersRepo.GetCharactersByUserID(user.ID)
	if err != nil {
		log.Printf("failed to find character: %v\n", err)
		return nil, err
	}

	return characters, nil
}

func (r *CharactersResolver) GetAllCharacters(params graphql.ResolveParams) (interface{}, error) {
	characters, err := r.CharactersRepo.GetAllCharacters()
	if err != nil {
		log.Printf("failed to find characters: %v\n", err)
		return nil, err
	}

	return characters, nil
}

func (r *CharactersResolver) UpdateCharacter(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if name, ok := params.Args["name"].(string); ok {
		character.Name = name
	}

	if tags, ok := params.Args["tags"].([]interface{}); ok {
		character.Tags = convertListToSlice(tags)
	}

	updateResult, err := r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return updateResult, nil
}

func (r *CharactersResolver) DeleteCharacter(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	deleteResult, err := r.CharactersRepo.DeleteCharacter(objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete character: %v", err)
	}

	return deleteResult, nil
}

func (r *CharactersResolver) ResetCharacter(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	character.Tags = []string{}
	character.TotalFocusedTime = 0
	character.CustomMetrics = []coredb.CustomMetric{}

	updateResult, err := r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to delete character: %v", err)
	}

	return updateResult, nil
}
