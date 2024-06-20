package characters

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createCharacter(params graphql.ResolveParams) (interface{}, error) {
	name := params.Args["name"].(string)
	var tags []string
	tagsInterface := params.Args["tags"].([]interface{})
	tags = make([]string, len(tagsInterface))
	for i, tag := range tagsInterface {
		tags[i] = tag.(string)
	}
	totalFocusTime := params.Args["totalFocusTime"].(int)

	customMetricsInput := params.Args["customMetrics"].([]interface{})
	var customMetricDatas []coredb.CustomMetric

	for _, cm := range customMetricsInput {
		cmMap, _ := cm.(map[string]interface{})
		metricID, _ := cmMap["id"].(primitive.ObjectID)
		metricName, _ := cmMap["name"].(string)
		metricDescription, _ := cmMap["description"].(string)
		metricTime, _ := cmMap["time"].(int)

		metricStyle, _ := cmMap["style"].(interface{})
		styleMap, _ := metricStyle.(map[string]interface{})
		styleColor, _ := styleMap["color"].(string)
		styleIcon, _ := styleMap["icon"].(string)
		styleData := coredb.StyleType{
			Color: styleColor,
			Icon:  styleIcon,
		}

		metricProperties, _ := cmMap["properties"].([]interface{})

		var propertiesData []coredb.MetricProperty
		for _, prop := range metricProperties {
			propMap, _ := prop.(map[string]interface{})
			propName, _ := propMap["name"].(string)
			propType, _ := propMap["type"].(string)
			propValue, _ := propMap["value"].(string)
			propUnit, _ := propMap["unit"].(string)
			var newValue interface{}
			switch propType {
			case "Text":
				newValue = propValue
			case "Number":
				newValue, _ = strconv.ParseFloat(propValue, 64)
			default:
				newValue = propValue
			}

			property := coredb.MetricProperty{
				Name:  propName,
				Type:  propType,
				Value: newValue,
				Unit:  propUnit,
			}
			propertiesData = append(propertiesData, property)
		}

		customMetric := coredb.CustomMetric{
			ID:          metricID,
			Name:        metricName,
			Description: metricDescription,
			Time:        int32(metricTime),
			Style:       styleData,
			Properties:  propertiesData,
		}
		customMetricDatas = append(customMetricDatas, customMetric)
	}

	character := coredb.Character{
		ID:               primitive.NewObjectID(),
		Name:             name,
		Tags:             tags,
		TotalFocusedTime: int32(totalFocusTime),
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
			metricName, _ := cmMap["name"].(string)
			metricDescription, _ := cmMap["description"].(string)
			metricTime, _ := cmMap["time"].(int)
			metricStyle, _ := cmMap["style"].(interface{})
			styleMap, _ := metricStyle.(map[string]interface{})
			styleColor, _ := styleMap["color"].(string)
			styleIcon, _ := styleMap["icon"].(string)
			styleData := coredb.StyleType{
				Color: styleColor,
				Icon:  styleIcon,
			}

			metricProperties, _ := cmMap["properties"].([]interface{})

			var propertiesData []coredb.MetricProperty
			for _, prop := range metricProperties {
				propMap, _ := prop.(map[string]interface{})
				propName, _ := propMap["name"].(string)
				propType, _ := propMap["type"].(string)
				propValue, _ := propMap["value"].(string)
				propUnit, _ := propMap["unit"].(string)
				var newValue interface{}
				switch propType {
				case "Text":
					newValue = propValue
				case "Number":
					newValue, _ = strconv.ParseFloat(propValue, 64)
				default:
					newValue = propValue
				}

				property := coredb.MetricProperty{
					Name:  propName,
					Type:  propType,
					Value: newValue,
					Unit:  propUnit,
				}
				propertiesData = append(propertiesData, property)
			}

			customMetric := coredb.CustomMetric{
				ID:          metricID,
				Name:        metricName,
				Description: metricDescription,
				Time:        int32(metricTime),
				Style:       styleData,
				Properties:  propertiesData,
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

func resetCharacter(params graphql.ResolveParams) (interface{}, error) {
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

	character.Tags = []string{}
	character.TotalFocusedTime = 0
	character.CustomMetrics = []coredb.CustomMetric{}
	update := bson.M{"$set": character}

	_, err = db.GetCharactersCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to delete character: %v", err)
	}

	return true, nil
}
