package character

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (db *Database) Connect() error {
	clientOptions := options.Client().ApplyURI("mongodb+srv://gotosleep:qAgVbxBQ03lkrQYx@tenk-hours-dev.bbehsco.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	log.Println("Connected to MongoDB!")

	db.client = client
	db.collection = client.Database("TenK-Hours-Dev").Collection("character")

	return nil
}

func (db *Database) Disconnect() {
	if db.client != nil {
		if err := db.client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v\n", err)
		}
	}
}

func (db *Database) InsertCharacter(character CharacterData) error {
	_, err := db.collection.InsertOne(context.Background(), character)
	return err
}

func (db *Database) UpdateCharacter(id string, character CharacterData) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": character}

	log.Printf("Updating character with ID: %s\n", id)

	result, err := db.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	log.Printf("Matched %d documents and modified %d documents\n", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (db *Database) GetCharacterByID(id string) (*CharacterData, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var character CharacterData
	filter := bson.M{"_id": objectID}
	err = db.collection.FindOne(context.Background(), filter).Decode(&character)
	if err != nil {
		return nil, err
	}
	return &character, nil
}

func (db *Database) GetAllCharacters() ([]*CharacterData, error) {
	var characters []*CharacterData
	cursor, err := db.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var character CharacterData
		if err := cursor.Decode(&character); err != nil {
			return nil, err
		}
		characters = append(characters, &character)
	}
	return characters, nil
}

func (db *Database) DeleteCharacter(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectID}
	result, err := db.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	log.Printf("Deleted %d documents\n", result.DeletedCount)

	return nil
}
