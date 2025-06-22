package mongorepo

import (
	"context"

	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharacterRepo struct {
	*mongodb.BaseRepo[entity.Character, mongomodel.Character]
}

func NewCharacterRepo(db *mongo.Database) *CharacterRepo {
	return &CharacterRepo{mongodb.NewBaseRepo[entity.Character, mongomodel.Character](
		db.Collection(mongodb.CharactersCollection),
		true,
	)}
}

func (r *CharacterRepo) GetCharactersByProfileID(ctx context.Context, profileID string) ([]entity.Character, error) {
	return r.FindMany(ctx, bson.M{"profile_id": mongodb.ToObjectID(profileID)})
}

func (r *CharacterRepo) CountCharactersByProfileID(ctx context.Context, profileID string) (int, error) {
	return r.Count(ctx, bson.M{"profile_id": mongodb.ToObjectID(profileID)})
}

func (r *CharacterRepo) GetAllCharacters(ctx context.Context) ([]entity.Character, error) {
	return r.FindMany(ctx, bson.M{})
}

func (r *CharacterRepo) DeleteCharacter(ctx context.Context, id string) (*entity.Character, error) {
	return r.FindAndDeleteByID(ctx, id)
}

func (r *CharacterRepo) DeleteCharactersByProfileID(ctx context.Context, profileID string) error {
	return r.DeleteMany(ctx, bson.M{"profile_id": mongodb.ToObjectID(profileID)})
}

func (r *CharacterRepo) Exist(ctx context.Context, profileID, characterID string) error {
	return r.Exists(ctx, bson.M{
		"_id":        mongodb.ToObjectID(characterID),
		"profile_id": mongodb.ToObjectID(profileID),
	})
}
