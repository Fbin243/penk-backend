package characters

import (
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *CharactersResolver) CreateCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object id: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.UserID != user.ID {
		return nil, auth.ErrorPermissionDenied
	}

	if len(character.CustomMetrics) >= int(character.LimitedMetricNumber) {
		return nil, fmt.Errorf("custom metric creation limit reached")
	}

	customMetric := coredb.CustomMetric{
		ID:                    primitive.NewObjectID(),
		Description:           "",
		Time:                  0,
		Style:                 coredb.MetricStyle{},
		Properties:            []coredb.MetricProperty{},
		LimitedPropertyNumber: 2,
	}

	input := params.Args["input"].(map[string]interface{})
	if name, ok := input["name"].(string); ok {
		customMetric.Name = name
	}

	if description, ok := input["description"].(string); ok {
		customMetric.Description = description
	}

	if metricStyle, ok := input["style"].(map[string]interface{}); ok {
		if color, ok := metricStyle["color"].(string); ok {
			customMetric.Style.Color = color
		}

		if icon, ok := metricStyle["icon"].(string); ok {
			customMetric.Style.Icon = icon
		}
	}

	err = ValidateCustomMetric(customMetric)
	if err != nil {
		return nil, err
	}

	createdCustomMetric, err := r.CharactersRepo.CreateCustomMetric(character.ID, &customMetric)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom metric: %v", err)
	}

	return *createdCustomMetric, nil
}

func (r *CharactersResolver) UpdateCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	metricID := params.Args["id"].(string)
	metricOID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.UserID != user.ID {
		return nil, auth.ErrorPermissionDenied
	}

	found := false
	updatedMetric := coredb.CustomMetric{}
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricOID {
			input := params.Args["input"].(map[string]interface{})
			if name, ok := input["name"].(string); ok {
				cm.Name = name
			}

			if description, ok := input["description"].(string); ok {
				cm.Description = description
			}

			if style, ok := input["style"].(map[string]interface{}); ok {
				if color, ok := style["color"].(string); ok {
					cm.Style.Color = color
				}

				if icon, ok := style["icon"].(string); ok {
					cm.Style.Icon = icon
				}
			}

			err = ValidateCustomMetric(cm)
			if err != nil {
				return nil, err
			}

			character.CustomMetrics[i] = cm
			updatedMetric = cm
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("custom metric does not belong to the character")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update custom metric: %v", err)
	}

	return updatedMetric, nil
}

func (r *CharactersResolver) DeleteCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character id: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.UserID != user.ID {
		return nil, auth.ErrorPermissionDenied
	}

	metricID := params.Args["id"].(string)
	metricOID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric id: %v", err)
	}

	found := false
	for _, cm := range character.CustomMetrics {
		if cm.ID == metricOID {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("custom metric does not belong to the character")
	}

	deletedCustomMetric, err := r.CharactersRepo.DeleteCustomMetric(characterOID, metricOID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete metric: %v", err)
	}

	return *deletedCustomMetric, nil
}

func (r *CharactersResolver) ResetCustomMetric(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	if characterOID != user.ID {
		return nil, auth.ErrorPermissionDenied
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	metricID := params.Args["id"].(string)
	metricOID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	found := false
	resetMetric := coredb.CustomMetric{}
	for i, metric := range character.CustomMetrics {
		if metric.ID == metricOID {
			metric.Description = ""
			metric.Time = 0
			metric.Style = coredb.MetricStyle{}
			metric.Properties = []coredb.MetricProperty{}

			character.CustomMetrics[i] = metric
			resetMetric = metric
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("custom metric does not belong to the character")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to reset custom metric: %v", err)
	}

	return resetMetric, nil
}

func (r *CharactersResolver) CreateMetricProperty(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.UserID != user.ID {
		return nil, auth.ErrorPermissionDenied
	}

	metricID := params.Args["metricID"].(string)
	metricOID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	metricProperty := coredb.MetricProperty{
		ID: primitive.NewObjectID(),
	}

	input := params.Args["input"].(map[string]interface{})
	if propName, ok := input["name"].(string); ok {
		metricProperty.Name = propName
	}

	if propType, ok := input["type"].(string); ok {
		metricProperty.Type = propType
	}

	if propValue, ok := input["value"].(string); ok {
		metricProperty.Value = castType(metricProperty.Type, propValue)
	}

	if propUnit, ok := input["unit"].(string); ok {
		metricProperty.Unit = propUnit
	}

	found := false
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricOID {
			if len(character.CustomMetrics[i].Properties) >= int(character.CustomMetrics[i].LimitedPropertyNumber) {
				return nil, fmt.Errorf("metric properties creation limit reached")
			}

			err = ValidateMetricProperty(metricProperty)
			if err != nil {
				return nil, err
			}

			character.CustomMetrics[i].Properties = append(character.CustomMetrics[i].Properties, metricProperty)
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("custom metric does not belong to the character")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric property: %v", err)
	}

	return metricProperty, nil
}

func (r *CharactersResolver) UpdateMetricProperty(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	metricID := params.Args["metricID"].(string)
	metricOID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	metricPropID := params.Args["id"].(string)
	metricPropOID, err := primitive.ObjectIDFromHex(metricPropID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric property ID: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.UserID != user.ID {
		return nil, auth.ErrorPermissionDenied
	}

	foundForMetric := false
	foundForProperty := false
	updatedProperty := coredb.MetricProperty{}
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricOID {
			for j, prop := range character.CustomMetrics[i].Properties {
				if prop.ID == metricPropOID {
					input := params.Args["input"].(map[string]interface{})
					if propName, ok := input["name"].(string); ok {
						prop.Name = propName
					}

					if propType, ok := input["type"].(string); ok {
						prop.Type = propType
					}

					if propValue, ok := input["value"].(string); ok {
						prop.Value = castType(prop.Type, propValue)
					}

					if propUnit, ok := input["unit"].(string); ok {
						prop.Unit = propUnit
					}

					err := ValidateMetricProperty(prop)
					if err != nil {
						return nil, err
					}

					character.CustomMetrics[i].Properties[j] = prop
					updatedProperty = prop
					foundForProperty = true
					break
				}
			}

			foundForMetric = true
			break
		}
	}

	if !foundForMetric {
		return nil, fmt.Errorf("custom metric does not belong to the character")
	}

	if !foundForProperty {
		return nil, fmt.Errorf("metric property does not belong to the metric")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update metric property: %v", err)
	}

	return updatedProperty, nil
}

func (r *CharactersResolver) DeleteMetricProperty(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	metricID := params.Args["metricID"].(string)
	metricOID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	metricPropID := params.Args["id"].(string)
	metricPropOID, err := primitive.ObjectIDFromHex(metricPropID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric property ID: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.UserID != user.ID {
		return nil, auth.ErrorPermissionDenied
	}

	foundForMetric := false
	foundForProperty := false
	deletedMetricProperty := coredb.MetricProperty{}
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricOID {
			for j, prop := range character.CustomMetrics[i].Properties {
				if prop.ID == metricPropOID {
					deletedMetricProperty = prop
					character.CustomMetrics[i].Properties = append(character.CustomMetrics[i].Properties[:j], character.CustomMetrics[i].Properties[j+1:]...)
					foundForProperty = true
					break
				}
			}

			foundForMetric = true
			break
		}
	}

	if !foundForMetric {
		return nil, fmt.Errorf("custom metric does not belong to the character")
	}

	if !foundForProperty {
		return nil, fmt.Errorf("metric property does not belong to the metric")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to delete metric property: %v", err)
	}

	return deletedMetricProperty, nil
}
