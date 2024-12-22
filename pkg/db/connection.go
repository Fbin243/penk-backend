package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ProfileCollection            = "profiles"
	CharacterCollection          = "characters"
	GoalsCollection              = "goals"
	TimeTrackingsCollection      = "time_trackings"
	SnapshotsCollection          = "snapshots"
	CapturedRecordsCollection    = "captured_records"
	FishCollection               = "fish"
	RodsCollection               = "rods"
	TemplatesCollection          = "templates"
	TemplateCategoriesCollection = "template_categories"
	FindOneAndUpdateOptions      = options.FindOneAndUpdate().SetReturnDocument(options.After)
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

func InitDBManagerFromEnv(dbName string) *DatabaseManager {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	mongoAddress := os.Getenv("MONGO_ADDRESS")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoDatabase := os.Getenv("MONGO_DATABASE_NAME")
	if dbName != "" {
		mongoDatabase = dbName
	}

	connectionURI := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/%s",
		mongoUser,
		mongoPassword,
		mongoAddress,
		mongoDatabase,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI).SetMaxConnIdleTime(5*60*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
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

	db.CreateCollection(ctx, CapturedRecordsCollection,
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
	if dbManager == nil {
		dbManager = InitDBManagerFromEnv("")
	}

	return dbManager
}
