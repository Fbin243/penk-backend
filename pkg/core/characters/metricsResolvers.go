package characters

import (
	"fmt"

	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetricsResolver struct {
	CharactersRepo *coredb.CharactersRepo
}

func NewMetricsResolver() *MetricsResolver {
	return &MetricsResolver{
		CharactersRepo: coredb.NewCharactersRepo(),
	}
}

func (r *CharactersResolver) CreateCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	characterID := params.Args["characterID"].(string)

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object id: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if len(character.CustomMetrics) >= int(character.LimitedMetricNumber) {
		fmt.Println(len(character.CustomMetrics))
		fmt.Println(character.LimitedMetricNumber)

		return nil, fmt.Errorf("custom metric creation limit reached")
	}

	name := params.Args["name"].(string)
	description := params.Args["description"].(string)
	style := params.Args["style"].(map[string]interface{})

	styleData := coredb.MetricStyle{
		Color: style["color"].(string),
		Icon:  style["icon"].(string),
	}

	customMetric := coredb.CustomMetric{
		ID:                    primitive.NewObjectID(),
		Name:                  name,
		Description:           description,
		Time:                  0,
		Style:                 styleData,
		Properties:            []coredb.MetricProperty{},
		LimitedPropertyNumber: 2,
	}

	result, err := r.CharactersRepo.AddCustomMetric(character.ID, customMetric)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom metric: %v", err)
	}

	return result, nil
}

func (r *CharactersResolver) UpdateCustomMetric(params graphql.ResolveParams) (interface{}, error) {
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

	character, err := r.CharactersRepo.GetCharacterByID(characterObjectID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	found := false
	for _, cm := range character.CustomMetrics {
		if cm.ID == metricObjectID {
			if name, ok := params.Args["name"].(string); ok {
				cm.Name = name
			}
			if description, ok := params.Args["description"].(string); ok {
				cm.Description = description
			}
			if style, ok := params.Args["style"].(map[string]interface{}); ok {
				cm.Style = coredb.MetricStyle{
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

				if len(propertiesData) > int(cm.LimitedPropertyNumber) {
					return nil, fmt.Errorf("custom metric properties creation limit reached")
				}
				cm.Properties = propertiesData
			}
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("custom metric not found")
	}

	updateResult, err := r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update custom metric: %v", err)
	}

	return updateResult, nil
}

func (r *CharactersResolver) DeleteCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	metricID := params.Args["id"].(string)
	characterID := params.Args["characterID"].(string)

	objectID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object id: %v", err)
	}

	_, err = r.CharactersRepo.GetCharacterByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	metricObjectID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	result, err := r.CharactersRepo.DeleteCustomMetric(objectID, metricObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete metric: %v", err)
	}

	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("metric not found or already deleted")
	}

	return true, nil
}

func (r *CharactersResolver) ResetCustomMetric(params graphql.ResolveParams) (interface{}, error) {
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

	character, err := r.CharactersRepo.GetCharacterByID(characterObjectID)
	if err != nil {
		return false, fmt.Errorf("character not found: %v", err)
	}

	found := false
	for _, metric := range character.CustomMetrics {
		if metric.ID == metricObjectID {
			metric.Description = ""
			metric.Time = 0
			metric.Style = coredb.MetricStyle{}
			metric.Properties = []coredb.MetricProperty{}
			found = true
			break
		}
	}

	if !found {
		return false, fmt.Errorf("custom metric not found")
	}

	updateResult, err := r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return false, fmt.Errorf("failed to reset custom metric: %v", err)
	}

	return updateResult, nil
}
