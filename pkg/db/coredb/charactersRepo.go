package coredb

import (
	"context"
	"log"
	"time"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharactersRepo struct {
	*mongo.Collection
}

func NewCharactersRepo(mongodb *mongo.Database) *CharactersRepo {
	return &CharactersRepo{mongodb.Collection(db.CharacterCollection)}
}

func (r *CharactersRepo) GetCharacterByID(id primitive.ObjectID) (*Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := Character{}
	err := r.FindOne(ctx, bson.M{"_id": id}).Decode(&character)

	return &character, err
}

func (r *CharactersRepo) GetCharactersByUserID(userID primitive.ObjectID) ([]Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var characters []Character
	err = cursor.All(ctx, &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (r *CharactersRepo) GetAllCharacters() ([]Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var characters []Character
	cursor, err := r.Find(ctx, primitive.M{})
	if err != nil {
		log.Printf("failed to fetch characters: %v\n", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &characters)
	if err != nil {
		log.Printf("failed to decode characters: %v\n", err)
		return nil, err
	}

	return characters, nil
}

func (r *CharactersRepo) CreateCharacter(character *Character) (*Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, character)

	return character, err
}

func (r *CharactersRepo) UpdateCharacter(character *Character) (*Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.FindOneAndUpdate(ctx, bson.M{"_id": character.ID}, bson.M{"$set": character}, db.FindOneAndUpdateOptions).Decode(character)

	return character, err
}

func (r *CharactersRepo) DeleteCharacter(id primitive.ObjectID) (*Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := &Character{}
	err := r.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(character)

	return character, err
}

func (r *CharactersRepo) CreateCustomMetric(characterID primitive.ObjectID, metric *CustomMetric) (*CustomMetric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.UpdateOne(ctx, bson.M{"_id": characterID}, bson.M{"$push": bson.M{"custom_metrics": *metric}})

	return metric, err
}

func (r *CharactersRepo) DeleteCustomMetric(characterID primitive.ObjectID, metricID primitive.ObjectID) (*CustomMetric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metric := &CustomMetric{}
	err := r.FindOneAndUpdate(ctx, bson.M{"_id": characterID}, bson.M{
		"$pull": bson.M{
			"custom_metrics": bson.M{"_id": metricID},
		},
	}).Decode(metric)

	return metric, err
}
