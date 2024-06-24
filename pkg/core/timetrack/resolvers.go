package timetrack

import (
	"context"
	"fmt"
	"log"
	"time"

	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createTimeTracking(params graphql.ResolveParams) (interface{}, error) {
	characterID := params.Args["characterID"].(string)
	startTime := time.Unix(int64(params.Args["startTime"].(int64)), 0)
	endTime := time.Unix(int64(params.Args["endTime"].(int)), 0)

	timeTracking := coredb.TimeTracking{
		ID:          primitive.NewObjectID(),
		CharacterID: characterID,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.GetTimeTrackingsCollection().InsertOne(ctx, timeTracking)
	if err != nil {
		log.Printf("Failed to insert time tracking: %v\n", err)
		return nil, err
	}

	err = updateCharacterTotalFocusTime(timeTracking.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("Failed to update character's total focus time: %v", err)
	}

	return timeTracking, nil
}

func getTimeTrackingByID(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	timeTracking := coredb.TimeTracking{}
	filter := bson.M{"_id": objectID}

	err = db.GetTimeTrackingsCollection().FindOne(ctx, filter).Decode(&timeTracking)
	if err != nil {
		log.Printf("Failed to find time tracking: %v\n", err)
		return nil, err
	}

	return timeTracking, nil
}

func getAllTimeTrackings(params graphql.ResolveParams) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var timeTrackings []coredb.TimeTracking
	cursor, err := db.GetTimeTrackingsCollection().Find(ctx, primitive.M{})
	if err != nil {
		log.Printf("Failed to fetch time trackings: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &timeTrackings)
	if err != nil {
		log.Printf("Failed to decode time trackings: %v\n", err)
		return nil, err
	}

	return timeTrackings, nil
}

func getTimeTrackingsByCharacterID(params graphql.ResolveParams) (interface{}, error) {
	characterID := params.Args["characterID"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var timeTrackings []coredb.TimeTracking
	filter := bson.M{"characterID": characterID}

	cursor, err := db.GetTimeTrackingsCollection().Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to fetch time trackings: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &timeTrackings)
	if err != nil {
		log.Printf("Failed to decode time trackings: %v\n", err)
		return nil, err
	}

	return timeTrackings, nil
}

func updateTimeTracking(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	timeTracking := coredb.TimeTracking{}
	filter := bson.M{"_id": objectID}

	err = db.GetTimeTrackingsCollection().FindOne(ctx, filter).Decode(&timeTracking)
	if err != nil {
		return nil, fmt.Errorf("Time tracking not found: %v", err)
	}

	if characterID, ok := params.Args["characterID"].(string); ok {
		timeTracking.CharacterID = characterID
	}

	if startTime, ok := params.Args["startTime"].(int64); ok {
		timeTracking.StartTime = time.Unix(startTime, 0)
	}

	if endTime, ok := params.Args["endTime"].(int64); ok {
		timeTracking.EndTime = time.Unix(endTime, 0) // Convert int64 to time.Time
	}

	update := bson.M{"$set": timeTracking}

	_, err = db.GetTimeTrackingsCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("Failed to update time tracking: %v", err)
	}

    // Update Character's focus time
    if customMetricsID, ok := params.Args["customMetricsID"].(string); ok && customMetricsID != "" {
        // Update specific custom metric
        customMetricObjectID, _ := primitive.ObjectIDFromHex(customMetricsID)
        filter := bson.M{
            "_id":           timeTracking.CharacterID,
            "customMetrics._id": customMetricObjectID,
        }

        // Calculate duration and update the specific custom metric's value
        duration := int64(timeTracking.EndTime.Sub(timeTracking.StartTime).Seconds())
        update := bson.M{"$inc": bson.M{"customMetrics.$.value": duration}}

        _, err = db.GetCharactersCollection().UpdateOne(ctx, filter, update)
        if err != nil {
            return nil, fmt.Errorf("failed to update custom metric: %v", err)
        }

        // Update total focus time as well (after updating custom metric)
        err = updateCharacterTotalFocusTime(timeTracking.CharacterID)
        if err != nil {
            return nil, fmt.Errorf("failed to update character's total focus time: %v", err)
        }

    } else {
        // Update total focus time
        err = updateCharacterTotalFocusTime(timeTracking.CharacterID)
        if err != nil {
            return nil, fmt.Errorf("failed to update character's total focus time: %v", err)
        }
        
    }

	err = updateCharacterTotalFocusTime(timeTracking.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("Failed to update character's total focus time: %v", err)
	}

	return timeTracking, nil
}

func deleteTimeTracking(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	timeTracking := coredb.TimeTracking{}
	err = db.GetTimeTrackingsCollection().FindOne(ctx, filter).Decode(&timeTracking)
	if err != nil {
		return nil, fmt.Errorf("Failed to find time tracking: %v", err)
	}

	_, err = db.GetTimeTrackingsCollection().DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to delete time tracking: %v", err)
	}

	err = updateCharacterTotalFocusTime(timeTracking.CharacterID)
	if err != nil {
		// Handle error, maybe log and potentially revert the deletion if possible
		return nil, fmt.Errorf("Failed to update character's total focus time after deletion: %v", err)
	}

	return true, nil
}

func updateCharacterTotalFocusTime(characterID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(characterID)

	totalFocusTime, err := calculateTotalFocusTime(characterID)
	if err != nil {
		return fmt.Errorf("Failed to calculate total focus time: %v", err)
	}

	update := bson.M{"$set": bson.M{"totalFocusTime": totalFocusTime}}
	_, err = db.GetCharactersCollection().UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("Failed to update character's total focus time in database: %v", err)
	}

	return nil
}

func calculateTotalFocusTime(characterID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"characterID": characterID}

	cursor, err := db.GetTimeTrackingsCollection().Find(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("Failed to fetch time trackings for character: %v", err)
	}
	defer cursor.Close(ctx)

	var timeTrackings []coredb.TimeTracking
	err = cursor.All(ctx, &timeTrackings)
	if err != nil {
		return 0, fmt.Errorf("Failed to decode time trackings: %v", err)
	}

	var totalFocusTime int64
	for _, tt := range timeTrackings {
		duration := tt.EndTime.Sub(tt.StartTime)
		totalFocusTime += int64(duration.Seconds())
	}

	return totalFocusTime, nil
}
