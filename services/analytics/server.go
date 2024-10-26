package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"tenkhours/pkg/analytics"
	"tenkhours/pkg/auth"
	"tenkhours/pkg/cron"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/analyticsdb"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/sessions"
	"tenkhours/services/analytics/graph"
	"tenkhours/services/analytics/graph/model"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("TENK_ENV")
	if env == "" {
		env = "development"
	}

	if godotenv.Load(".env."+env) != nil {
		log.Fatal("Error loading .env." + env + " file")
	}

	fmt.Println("------------------Running in environment:", env)

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	// Init dependencies and perform DI manually
	db := db.GetDBManager().DB
	profilesRepo := coredb.NewProfilesRepo(db)
	charactersRepo := coredb.NewCharactersRepo(db)
	snapshotsRepo := analyticsdb.NewSnapshotRepo(db)
	capturedRecordsRepo := analyticsdb.NewCapturedRecordRepo(db)
	charactersHandler := analytics.NewCharactersHandler(snapshotsRepo, charactersRepo, profilesRepo, capturedRecordsRepo)
	redisClient := sessions.GetRedisClient()

	// Make a cron run daily for captured records
	cron := cron.NewCron()
	cron.RunDaily(func() {
		fmt.Println("Running cron job every day")

		// Scan Redis and save to DB
		for {
			var cursor uint64
			profileIds, cursor, err := redisClient.Scan(context.Background(), cursor, "*"+sessions.CapturedRecordKey+"*", 1000).Result()
			if err != nil {
				fmt.Println("Error scanning redis: ", err)
			}
			
			for _, profileId := range profileIds {
				// Get the captured record from redis and save it to DB
				capturedRecordsJSON, err := redisClient.HGetAll(context.Background(), profileId).Result()
				if err != nil {
					fmt.Println("Error getting characters from redis: ", err)
				}

				for _, capturedRecordJSON := range capturedRecordsJSON {
					var capturedRecord model.CapturedRecord
					// Decode the captured records json to struct
					err = json.Unmarshal([]byte(capturedRecordJSON), &capturedRecord)
					if err != nil {
						fmt.Println("Error unmarshalling captured record: ", err)
					}

					// Save the captured record to DB
					err = capturedRecordsRepo.CreateCapturedRecord(&capturedRecord)
					if err != nil {
						fmt.Println("Error saving captured record to DB: ", err)
					}
				}

				// Delete the captured records from redis
				_, err = redisClient.Del(context.Background(), profileId).Result()
				if err != nil {
					fmt.Println("Error deleting captured records from redis:", err)
				}
			}
			if cursor == 0 {
				// cron.Stop() // Stop the cron job for testing
				break
			}
		}
	})

	// Check authentication
	authMiddleware := auth.NewMiddleware(profilesRepo)
	app.Use(authMiddleware.CheckAuth)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CharactersHandler: charactersHandler,
		},
	}))

	// app.GET("/", func(c *gin.Context) {
	// 	playgroundHandler := playground.Handler("GraphQL playground", "/graphql")
	// 	playgroundHandler.ServeHTTP(c.Writer, c.Request)
	// })

	app.POST("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	port, found := os.LookupEnv("ANALYTICS_PORT")
	if !found {
		port = "8083"
	}
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	app.Run(":" + port)
}
