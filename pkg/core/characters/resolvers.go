package characters

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

func createCharacter(params graphql.ResolveParams) (interface{}, error) {
	userID := params.Args["userID"].(string)

	name := params.Args["name"].(string)

	var tags []string
	tagsInterface := params.Args["tags"].([]interface{})
	tags = make([]string, len(tagsInterface))
	for i, tag := range tagsInterface {
		tags[i] = tag.(string)
	}

	totalFocusTime := int32(params.Args["totalFocusTime"].(int))

	customMetricsInput := params.Args["customMetrics"].([]interface{})

	var customMetricDatas []coredb.CustomMetric

	for _, cm := range customMetricsInput {
		cmMap, _ := cm.(map[string]interface{})
		metricID, _ := cmMap["id"].(primitive.ObjectID)
		metricCharacterID, _ := cmMap["characterID"].(primitive.ObjectID)
		metricType, _ := cmMap["type"].(string)
		metricName, _ := cmMap["name"].(string)
		metricValue, _ := cmMap["value"].(string)

		customMetric := coredb.CustomMetric{
			ID:          metricID,
			CharacterID: metricCharacterID,
			Type:        metricType,
			Name:        metricName,
			Value:       metricValue,
		}
		customMetricDatas = append(customMetricDatas, customMetric)
	}

	character := coredb.Character{
		ID:               primitive.NewObjectID(),
		UserID:           userID,
		Name:             name,
		Tags:             tags,
		TotalFocusedTime: totalFocusTime,
		CustomMetrics:    customMetricDatas,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.GetCharactersCollection().InsertOne(ctx, character)
	if err != nil {
		log.Printf("Failed to insert character: %v\n", err)
		return nil, err
	}

	return character, nil
}

func getCharacterByID(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := coredb.Character{}
	filter := bson.M{"_id": objectID}

	err = db.GetCharactersCollection().FindOne(ctx, filter).Decode(&character)
	if err != nil {
		log.Printf("Failed to find character: %v\n", err)
		return nil, err
	}

	return character, nil
}

func getAllCharacters(params graphql.ResolveParams) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var characters []coredb.Character
	cursor, err := db.GetCharactersCollection().Find(ctx, primitive.M{})
	if err != nil {
		log.Printf("Failed to fetch characters: %v\n", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &characters)
	if err != nil {
		log.Printf("Failed to decode characters: %v\n", err)
		return nil, err
	}

	return characters, nil
}

func updateCharacter(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := coredb.Character{}
	filter := bson.M{"_id": objectID}

	err = db.GetCharactersCollection().FindOne(ctx, filter).Decode(&character)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if userID, ok := params.Args["userID"].(string); ok {
		character.UserID = userID
	}

	if name, ok := params.Args["name"].(string); ok {
		character.Name = name
	}

	if tags, ok := params.Args["tags"].([]string); ok {
		character.Tags = tags
	}

	if totalFocusTime, ok := params.Args["totalFocusTime"].(int32); ok {
		character.TotalFocusedTime = totalFocusTime
	}

	if customMetricsInput, ok := params.Args["customMetrics"].([]interface{}); ok {
		var customMetricDatas []coredb.CustomMetric

		for _, cm := range customMetricsInput {
			cmMap, _ := cm.(map[string]interface{})
			metricID, _ := cmMap["id"].(primitive.ObjectID)
			metricCharacterID, _ := cmMap["characterID"].(primitive.ObjectID)
			metricType, _ := cmMap["type"].(string)
			metricName, _ := cmMap["name"].(string)
			metricValue, _ := cmMap["value"].(string)

			customMetric := coredb.CustomMetric{
				ID:          metricID,
				CharacterID: metricCharacterID,
				Type:        metricType,
				Name:        metricName,
				Value:       metricValue,
			}
			customMetricDatas = append(customMetricDatas, customMetric)
		}
		character.CustomMetrics = customMetricDatas
	}

	update := bson.M{"$set": character}

	_, err = db.GetCharactersCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return character, nil
}

func deleteCharacter(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.GetCharactersCollection().DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete character: %v", err)
	}

	return true, nil
}
