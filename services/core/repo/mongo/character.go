package mongorepo

import (
	"context"
	"fmt"
	"time"

	"tenkhours/services/core/entity"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharacterRepo struct {
	*mongodb.BaseRepo[entity.Character, Character]
}

func NewCharacterRepo(db *mongo.Database) *CharacterRepo {
	return &CharacterRepo{mongodb.NewBaseRepo(
		db.Collection(mongodb.CharactersCollection),
		&mongodb.Mapper[entity.Character, Character]{},
	)}
}

func (r *CharacterRepo) GetCharactersByProfileID(ctx context.Context, profileID string) ([]entity.Character, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"profile_id": mongodb.ToObjectID(profileID)})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var characters []entity.Character
	err = cursor.All(ctx, &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (r *CharacterRepo) CountCharactersByProfileID(ctx context.Context, profileID string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.CountDocuments(ctx, bson.M{"profile_id": mongodb.ToObjectID(profileID)})
}

func (r *CharacterRepo) GetAllCharacters(ctx context.Context) ([]entity.Character, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var characters []entity.Character
	cursor, err := r.Find(ctx, primitive.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (r *CharacterRepo) DeleteCharacter(ctx context.Context, id string) (*entity.Character, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	character := &entity.Character{}
	err := r.FindOneAndDelete(ctx, bson.M{"_id": mongodb.ToObjectID(id)}).Decode(character)

	return character, err
}

func (r *CharacterRepo) DeleteCharactersByProfileID(ctx context.Context, profileID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"profile_id": mongodb.ToObjectID(profileID)})

	return err
}
