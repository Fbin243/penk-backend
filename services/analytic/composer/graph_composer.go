package composer

import (
	"context"
	"encoding/json"
	"fmt"

	"tenkhours/pkg/cron"
	"tenkhours/services/analytic/entity"
	"tenkhours/services/analytic/transport/graph"

	rdb "tenkhours/pkg/db/redis"
)

func ComposeGraphQLResolver() *graph.Resolver {
	return &graph.Resolver{
		AnalyticBusiness: GetComposer().AnalyticBiz,
	}
}

func ComposeCronJobs() {
	// Make a cron run daily for captured records
	cron := cron.NewCron()
	cron.RunDaily(func() {
		fmt.Println("Running cron job every day")

		// Scan Redis and save to DB
		for {
			var cursor uint64
			profileIds, cursor, err := rdb.GetRedisClient().Scan(context.Background(), cursor, "*"+rdb.CapturedRecordKey+"*", 1000).Result()
			if err != nil {
				fmt.Println("Error scanning redis: ", err)
			}

			for _, profileId := range profileIds {
				// Get the captured record from redis and save it to DB
				capturedRecordsJSON, err := rdb.GetRedisClient().HGetAll(context.Background(), profileId).Result()
				if err != nil {
					fmt.Println("Error getting characters from redis: ", err)
				}

				for _, capturedRecordJSON := range capturedRecordsJSON {
					var capturedRecord entity.CapturedRecord
					// Decode the captured records json to struct
					err = json.Unmarshal([]byte(capturedRecordJSON), &capturedRecord)
					if err != nil {
						fmt.Println("Error unmarshalling captured record: ", err)
					}

					// Save the captured record to DB
					err = GetComposer().CapturedRecordRepo.CreateCapturedRecord(context.Background(), &capturedRecord)
					if err != nil {
						fmt.Println("Error saving captured record to DB: ", err)
					}
				}

				// Delete the captured records from redis
				_, err = rdb.GetRedisClient().Del(context.Background(), profileId).Result()
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
