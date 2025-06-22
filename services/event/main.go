package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	// Load environment variables from .env file
	env := os.Getenv("TENK_ENV")
	if env == "" {
		env = "test"
	}

	if godotenv.Load(".env."+env) != nil {
		log.Printf("Error loading .env." + env + " file")
	}
}

func StartWatcher(ctx context.Context) {
	colls := []string{
		mongodb.CategoriesCollection,
		mongodb.MetricsCollection,
		mongodb.GoalsCollection,
		mongodb.TimeTrackingsCollecion,
		mongodb.HabitsCollection,
		mongodb.HabitLogsCollection,
		mongodb.TasksCollection,
		mongodb.TaskSessionsCollection,
	}

	for _, coll := range colls {
		go watchCollection(ctx, coll)
	}
}

func watchCollection(ctx context.Context, collectionName string) {
	db := mongodb.GetDBManager().DB

	coll := db.Collection(collectionName)
	eventColl := db.Collection(mongodb.StreamEventsCollection)

	stream, err := coll.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		log.Printf("error watching %s: %v", collectionName, err)
		return
	}
	defer stream.Close(ctx)

	log.Printf("Watching collection: %s", collectionName)

	for stream.Next(ctx) {
		var event bson.M
		if err := stream.Decode(&event); err != nil {
			log.Println("Failed to decode event:", err)
			continue
		}

		log.Printf("Event detected in %s: %v", collectionName, utils.PrettyJSON(event))

		_, err := eventColl.InsertOne(ctx, event)
		if err != nil {
			log.Println("Failed to write log:", err)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartWatcher(ctx)

	// Handle SIGINT to gracefully shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down watcher...")
	cancel()
}
