package composer

import (
	"context"
	"encoding/json"
	"fmt"
	"tenkhours/pkg/cron"
	"tenkhours/pkg/db"
	"tenkhours/services/analytics/business"
	"tenkhours/services/analytics/graph"
	"tenkhours/services/analytics/graph/model"
	analyticsRepo "tenkhours/services/analytics/repo"
	"tenkhours/services/core/repo"

	"github.com/go-redis/redis/v8"
)

func ComposeGraphQLResolver() *graph.Resolver {
	// Init dependencies and perform DI manually
	mongodb := db.GetDBManager().DB
	redisClient := db.GetRedisClient()
	charactersRepo := repo.NewCharactersRepo(mongodb)
	profilesRepo := repo.NewProfilesRepo(mongodb, redisClient)
	snapshotsRepo := analyticsRepo.NewSnapshotsRepo(mongodb)
	capturedRecordsRepo := analyticsRepo.NewCapturedRecordsRepo(mongodb)
	charactersBiz := business.NewCharactersBusiness(snapshotsRepo, charactersRepo, profilesRepo, capturedRecordsRepo, redisClient)
	ComposeCronJobs(redisClient, capturedRecordsRepo)

	return &graph.Resolver{
		CharactersBusiness: charactersBiz,
	}
}

func ComposeCronJobs(redisClient *redis.Client, capturedRecordsRepo *analyticsRepo.CapturedRecordsRepo) {
	// Make a cron run daily for captured records
	cron := cron.NewCron()
	cron.RunDaily(func() {
		fmt.Println("Running cron job every day")

		// Scan Redis and save to DB
		for {
			var cursor uint64
			profileIds, cursor, err := redisClient.Scan(context.Background(), cursor, "*"+db.CapturedRecordKey+"*", 1000).Result()
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
}
