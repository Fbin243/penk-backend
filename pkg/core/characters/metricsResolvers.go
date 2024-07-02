package characters

import (
	"fmt"

	"tenkhours/pkg/db/coredb"
	"tenkhours/test"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *CharactersResolver) CreateCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return false, fmt.Errorf("failed to get object id: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return false, fmt.Errorf("failed to get character: %v", err)
	}

	if len(character.CustomMetrics) >= int(character.LimitedMetricNumber) {
		return false, fmt.Errorf("custom metric creation limit reached")
	}

	name := params.Args["name"].(string)
	description, ok := params.Args["description"].(string)
	if !ok {
		description = ""
	}

	style, ok := params.Args["style"].(map[string]interface{})
	if !ok {
		style = map[string]interface{}{
			"color": "",
			"icon":  "",
		}
	}

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

	// --> JUST FOR TESTING
	mTest := test.NewTestManager()
	ctx := mTest.GetContext()
	ctx.IdCustomMetric = customMetric.ID.Hex()
	mTest.UpdateContext(ctx)

	_, err = r.CharactersRepo.AddCustomMetric(character.ID, customMetric)
	if err != nil {
		return false, fmt.Errorf("failed to create custom metric: %v", err)
	}

	return true, nil
}

func (r *CharactersResolver) UpdateCustomMetric(params graphql.ResolveParams) (interface{}, error) {
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
	for i, cm := range character.CustomMetrics {
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

			character.CustomMetrics[i] = cm
			found = true
			break
		}
	}

	if !found {
		return false, fmt.Errorf("custom metric not found")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return false, fmt.Errorf("failed to update custom metric: %v", err)
	}

	return true, nil
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
	for i, metric := range character.CustomMetrics {
		if metric.ID == metricObjectID {
			metric.Description = ""
			metric.Time = 0
			metric.Style = coredb.MetricStyle{}
			metric.Properties = []coredb.MetricProperty{}

			character.CustomMetrics[i] = metric
			found = true
			break
		}
	}

	if !found {
		return false, fmt.Errorf("custom metric not found")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return false, fmt.Errorf("failed to reset custom metric: %v", err)
	}

	return true, nil
}

func (r *CharactersResolver) CreateMetricProperty(params graphql.ResolveParams) (interface{}, error) {
	metricID := params.Args["metricID"].(string)
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

	propName := params.Args["name"].(string)
	propType := params.Args["type"].(string)
	propValue := params.Args["value"].(string)
	propUnit, ok := params.Args["unit"].(string)
	if !ok {
		propUnit = ""
	}

	found := false
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricObjectID {
			if len(character.CustomMetrics[i].Properties) >= int(character.CustomMetrics[i].LimitedPropertyNumber) {
				return false, fmt.Errorf("metric properties creation limit reached")
			}

			metricProperty := coredb.MetricProperty{
				ID:    primitive.NewObjectID(),
				Name:  propName,
				Type:  propType,
				Value: castType(propType, propValue),
				Unit:  propUnit,
			}

			character.CustomMetrics[i].Properties = append(character.CustomMetrics[i].Properties, metricProperty)
			found = true
			break
		}
	}

	if !found {
		return false, fmt.Errorf("custom metric not found")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return false, fmt.Errorf("failed to create metric property: %v", err)
	}

	return true, nil
}

func (r *CharactersResolver) UpdateMetricProperty(params graphql.ResolveParams) (interface{}, error) {
	metricPropID := params.Args["id"].(string)
	metricPropObjectID, err := primitive.ObjectIDFromHex(metricPropID)
	if err != nil {
		return false, fmt.Errorf("invalid metric property ID: %v", err)
	}

	metricID := params.Args["metricID"].(string)
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

	foundForMetric := false
	foundForProperty := false
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricObjectID {
			for j, prop := range character.CustomMetrics[i].Properties {
				if prop.ID == metricPropObjectID {
					if propName, ok := params.Args["name"].(string); ok {
						prop.Name = propName
					}

					if propType, ok := params.Args["type"].(string); ok {
						prop.Type = propType
					}

					if propValue, ok := params.Args["value"].(string); ok {
						prop.Value = castType(prop.Type, propValue)
					}

					if propUnit, ok := params.Args["unit"].(string); ok {
						prop.Unit = propUnit
					}

					character.CustomMetrics[i].Properties[j] = prop
					foundForMetric = true
					break
				}
			}
			foundForMetric = true
			break
		}
	}

	if !foundForMetric {
		return false, fmt.Errorf("custom metric not found")
	}

	if !foundForProperty {
		return false, fmt.Errorf("metric property not found")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return false, fmt.Errorf("failed to update metric property: %v", err)
	}

	return true, nil
}

func (r *CharactersResolver) DeleteMetricProperty(params graphql.ResolveParams) (interface{}, error) {
	metricPropID := params.Args["id"].(string)
	metricPropObjectID, err := primitive.ObjectIDFromHex(metricPropID)
	if err != nil {
		return false, fmt.Errorf("invalid metric property ID: %v", err)
	}

	metricID := params.Args["metricID"].(string)
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

	foundForMetric := false
	foundForProperty := false
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricObjectID {
			for j, prop := range character.CustomMetrics[i].Properties {
				if prop.ID == metricPropObjectID {
					character.CustomMetrics[i].Properties = append(character.CustomMetrics[i].Properties[:j], character.CustomMetrics[i].Properties[j+1:]...)
					foundForMetric = true
					break
				}
			}
			foundForMetric = true
			break
		}
	}

	if !foundForMetric {
		return false, fmt.Errorf("custom metric not found")
	}

	if !foundForProperty {
		return false, fmt.Errorf("metric property not found")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return false, fmt.Errorf("failed to remove metric property: %v", err)
	}

	return true, nil
}
