package characters

import (
	"context"
	"fmt"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"
	"time"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	characterID := params.Args["characterID"].(string)

	objectID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object id: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := coredb.Character{}
	filter := bson.M{"_id": objectID}

	err = db.GetCharactersCollection().FindOne(ctx, filter).Decode(&character)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if len(character.CustomMetrics) >= int(character.LimitedMetrics) {
		return nil, fmt.Errorf("custom metric creation limit reached")
	}

	name := params.Args["name"].(string)
	description := params.Args["description"].(string)
	style := params.Args["style"].(map[string]interface{})

	styleData := coredb.StyleType{
		Color: style["color"].(string),
		Icon:  style["icon"].(string),
	}

	customMetric := coredb.CustomMetric{
		ID:                primitive.NewObjectID(),
		Name:              name,
		Description:       description,
		Time:              0,
		Style:             styleData,
		Properties:        []coredb.MetricProperty{},
		LimitedProperties: 2,
	}

	character.CustomMetrics = append(character.CustomMetrics, customMetric)

	update := bson.M{"$set": character}

	_, err = db.GetCharactersCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom metric: %v", err)
	}

	return customMetric, nil
}

func updateCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	metricID := params.Args["id"].(string)
	metricObjectID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	characterID := params.Args["characterID"].(string)
	characterObjectID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := coredb.Character{}
	filter := bson.M{"_id": characterObjectID}

	err = db.GetCharactersCollection().FindOne(ctx, filter).Decode(&character)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	var updatedMetric coredb.CustomMetric
	found := false
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricObjectID {
			if name, ok := params.Args["name"].(string); ok {
				character.CustomMetrics[i].Name = name
			}
			if description, ok := params.Args["description"].(string); ok {
				character.CustomMetrics[i].Description = description
			}
			if style, ok := params.Args["style"].(map[string]interface{}); ok {
				character.CustomMetrics[i].Style = coredb.StyleType{
					Color: style["color"].(string),
					Icon:  style["icon"].(string),
				}
			}
			if properites, ok := params.Args["properties"].([]interface{}); ok {
				var propertiesData []coredb.MetricProperty
				for _, prop := range properites {
					propMap, _ := prop.(map[string]interface{})
					propName, _ := propMap["name"].(string)
					propType, _ := propMap["type"].(string)
					propValue, _ := propMap["value"].(string)
					propUnit, _ := propMap["unit"].(string)

					property := coredb.MetricProperty{
						Name:  propName,
						Type:  propType,
						Value: castType(propType, propValue),
						Unit:  propUnit,
					}

					propertiesData = append(propertiesData, property)
				}

				if len(propertiesData) > int(character.CustomMetrics[i].LimitedProperties) {
					return nil, fmt.Errorf("custom metric properties creation limit reached")
				}
				character.CustomMetrics[i].Properties = propertiesData
			}
			updatedMetric = character.CustomMetrics[i]
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("custom metric not found")
	}

	update := bson.M{"$set": character}

	_, err = db.GetCharactersCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update custom metric: %v", err)
	}
	return updatedMetric, nil

}

func deleteCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	metricID := params.Args["id"].(string)
	characterID := params.Args["characterID"].(string)

	objectID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object id: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := coredb.Character{}
	filter := bson.M{"_id": objectID}

	err = db.GetCharactersCollection().FindOne(ctx, filter).Decode(&character)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	metricObjectID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	update := bson.M{
		"$pull": bson.M{
			"custom_metrics": bson.M{"_id": metricObjectID},
		},
	}

	result, err := db.GetCharactersCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to delete metric: %v", err)
	}

	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("metric not found or already deleted")
	}

	return true, nil
}

func resetCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	metricID := params.Args["id"].(string)
	metricObjectID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return false, fmt.Errorf("invalid metric ID: %v", err)
	}

	characterID := params.Args["characterID"].(string)
	characterObjectID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return false, fmt.Errorf("invalid character ID: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := coredb.Character{}
	filter := bson.M{"_id": characterObjectID}

	err = db.GetCharactersCollection().FindOne(ctx, filter).Decode(&character)
	if err != nil {
		return false, fmt.Errorf("character not found: %v", err)
	}

	found := false
	for i, metric := range character.CustomMetrics {
		if metric.ID == metricObjectID {
			character.CustomMetrics[i].Description = ""
			character.CustomMetrics[i].Time = 0
			character.CustomMetrics[i].Style = coredb.StyleType{}
			character.CustomMetrics[i].Properties = []coredb.MetricProperty{}
			found = true
			break
		}
	}

	if !found {
		return false, fmt.Errorf("custom metric not found")
	}

	update := bson.M{"$set": character}
	_, err = db.GetCharactersCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return false, fmt.Errorf("failed to reset custom metric: %v", err)
	}

	return true, nil
}
