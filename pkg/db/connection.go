package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db                  *mongo.Database
	UserCollection      = "users"
	CharacterCollection = "character"
)

func GetDB() *mongo.Database {
	if db != nil {
		return db
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mongoAddress := os.Getenv("MONGO_ADDRESS")
	mongoDatabase := os.Getenv("MONGO_DATABASE_NAME")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	connectionURI := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/%s",
		mongoUser,
		mongoPassword,
		mongoAddress,
		mongoDatabase,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(mongoDatabase)

	// optional setup for db here
	usersCollection := db.Collection(UserCollection)
	_, _ = usersCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "firebase_uid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	return db
}

func GetUsersCollection() *mongo.Collection {
	return GetDB().Collection(UserCollection)
}

func GetCharactersCollection() *mongo.Collection {
	return GetDB().Collection(CharacterCollection)
}
