package timetrackings

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
	customMetricID, ok := params.Args["customMetricID"].(string)
	if !ok {
		customMetricID = ""
	}

	timeTracking := coredb.TimeTracking{
		ID:             primitive.NewObjectID(),
		CharacterID:    characterID,
		CustomMetricID: customMetricID,
		StartTime:      time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.GetTimeTrackingsCollection().InsertOne(ctx, timeTracking)
	if err != nil {
		log.Printf("failed to insert time tracking: %v\n", err)
		return nil, err
	}

	return timeTracking, nil
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
		return nil, fmt.Errorf("time tracking not found: %v", err)
	}

	endTime := time.Now()
	characterID, err := primitive.ObjectIDFromHex(timeTracking.CharacterID)
	if err != nil {
		return nil, err
	}

	character := coredb.Character{}

	filter = bson.M{"_id": characterID}
	err = db.GetCharactersCollection().FindOne(ctx, filter).Decode(&character)
	if err != nil {
		return nil, err
	}

	duration := endTime.Sub(timeTracking.StartTime).Seconds()

	character.TotalFocusedTime += int32(duration)
	if timeTracking.CustomMetricID != "" {
		for i, customMetric := range character.CustomMetrics {
			if customMetric.ID.Hex() == timeTracking.CustomMetricID {
				character.CustomMetrics[i].Time += int32(duration)
				break
			}
		}
	}

	timeTracking.EndTime = endTime
	update := bson.M{"$set": timeTracking}
	_, err = db.GetTimeTrackingsCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update time tracking: %v", err)
	}

	update = bson.M{"$set": character}
	_, err = db.GetCharactersCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update character tracking: %v", err)
	}

	return timeTracking, nil
}
