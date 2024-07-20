package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	UserCollection          = "users"
	CharacterCollection     = "character"
	TimeTrackingsCollection = "time_trackings"
	SnapshotsCollection     = "snapshots"
	FindOneAndUpdateOptions = options.FindOneAndUpdate().SetReturnDocument(options.After)
)

type DatabaseManager struct {
	DB     *mongo.Database
	Client *mongo.Client
}

func InitDBManagerFromURL(connectionURI string, dbName string) *DatabaseManager {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		log.Fatal(err)
	}

	return &DatabaseManager{
		DB:     client.Database(dbName),
		Client: client,
	}
}

func InitDBManagerFromEnv() *DatabaseManager {
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

	db := client.Database(mongoDatabase)

	// TODO: Create time series collection for snapshots (Temporarily)
	db.CreateCollection(ctx, SnapshotsCollection,
		options.CreateCollection().
			SetTimeSeriesOptions(
				options.TimeSeries().
					SetTimeField("timestamp").
					SetMetaField("metadata"),
			),
	)

	return &DatabaseManager{
		DB:     db,
		Client: client,
	}
}

var dbManager *DatabaseManager

func GetDBManager() *DatabaseManager {
	if dbManager != nil {
		return dbManager
	}

	return InitDBManagerFromEnv()
}
