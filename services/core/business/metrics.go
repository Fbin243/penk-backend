package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *CharactersHandler) CreateCustomMetric(ctx context.Context, characterID primitive.ObjectID, input model.CustomMetricInput) (*repo.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	if len(character.CustomMetrics) >= int(character.LimitedMetricNumber) {
		return nil, fmt.Errorf("custom metric creation limit reached")
	}

	customMetric := repo.CustomMetric{
		ID:                    primitive.NewObjectID(),
		Description:           "",
		Time:                  0,
		Style:                 repo.MetricStyle{},
		Properties:            []repo.MetricProperty{},
		LimitedPropertyNumber: 2,
	}

	if input.Name != nil {
		customMetric.Name = *input.Name
	}

	if input.Description != nil {
		customMetric.Description = *input.Description
	}

	if input.Style != nil {
		if input.Style.Color != nil {
			customMetric.Style.Color = *input.Style.Color
		}

		if input.Style.Icon != nil {
			customMetric.Style.Icon = *input.Style.Icon
		}
	}

	if input.Properties != nil {
		var properties []repo.MetricProperty
		for _, prop := range input.Properties {
			var metricProperty repo.MetricProperty
			metricProperty.ID = primitive.NewObjectID()

			if prop.Name != nil {
				metricProperty.Name = *prop.Name
			}
			if prop.Type != nil {
				metricProperty.Type = *prop.Type
			}
			if prop.Value != nil {
				metricProperty.Value = *prop.Value
			}
			if prop.Unit != nil {
				metricProperty.Unit = *prop.Unit
			}

			if len(properties) >= int(customMetric.LimitedPropertyNumber) {
				return nil, fmt.Errorf("metric properties creation limit reached")
			}

			properties = append(properties, metricProperty)
		}

		customMetric.Properties = properties
	}

	createdCustomMetric, err := r.CharactersRepo.CreateCustomMetric(character.ID, &customMetric)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom metric: %v", err)
	}

	return createdCustomMetric, nil
}

func (r *CharactersHandler) UpdateCustomMetric(ctx context.Context, metricID primitive.ObjectID, characterID primitive.ObjectID, input model.CustomMetricInput) (*repo.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	found := false
	updatedMetric := repo.CustomMetric{}
	for i, cm := range character.CustomMetrics {
		if cm.ID != metricID {
			continue
		}
		if input.Name != nil {
			cm.Name = *input.Name
		}
		if input.Description != nil {
			cm.Description = *input.Description
		}
		if input.Style != nil {
			if input.Style.Color != nil {
				cm.Style.Color = *input.Style.Color
			}

			if input.Style.Icon != nil {
				cm.Style.Icon = *input.Style.Icon
			}
		}

		if input.Properties != nil {
			var properties []repo.MetricProperty
			for _, prop := range input.Properties {
				var metricProperty repo.MetricProperty
				if prop.ID != nil {
					metricProperty.ID = *prop.ID
					propertyFound := false
					for _, p := range cm.Properties {
						if p.ID == *prop.ID {
							metricProperty = p
							propertyFound = true
							break
						}
					}

					if !propertyFound {
						return nil, fmt.Errorf("metric property does not belong to the metric")
					}

				} else {
					metricProperty.ID = primitive.NewObjectID()
				}

				if prop.Name != nil {
					metricProperty.Name = *prop.Name
				}
				if prop.Type != nil {
					metricProperty.Type = *prop.Type
				}
				if prop.Value != nil {
					metricProperty.Value = *prop.Value
				}
				if prop.Unit != nil {
					metricProperty.Unit = *prop.Unit
				}

				if len(properties) >= int(character.CustomMetrics[i].LimitedPropertyNumber) {
					return nil, fmt.Errorf("metric properties creation limit reached")
				}

				properties = append(properties, metricProperty)
			}
			cm.Properties = properties
		}

		character.CustomMetrics[i] = cm
		updatedMetric = cm
		found = true
		break
	}

	if !found {
		return nil, fmt.Errorf("custom metric does not belong to the character")
	}

	_, err = r.CharactersRepo.UpdateCharacter(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update custom metric: %v", err)
	}

	return &updatedMetric, nil
}

func (r *CharactersHandler) DeleteCustomMetric(ctx context.Context, metricID primitive.ObjectID, characterID primitive.ObjectID) (*repo.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	found := false
	for _, cm := range character.CustomMetrics {
		if cm.ID == metricID {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("custom metric does not belong to the character")
	}

	deletedCustomMetric, err := r.CharactersRepo.DeleteCustomMetric(characterID, metricID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete metric: %v", err)
	}

	return deletedCustomMetric, nil
}

func (r *CharactersHandler) ResetCustomMetric(ctx context.Context, metricID primitive.ObjectID, characterID primitive.ObjectID) (*repo.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	found := false
	resetMetric := repo.CustomMetric{}
	for i, metric := range character.CustomMetrics {
		if metric.ID == metricID {
			metric.Description = ""
			metric.Time = 0
			metric.Style = repo.MetricStyle{}
			metric.Properties = []repo.MetricProperty{}

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

	return &resetMetric, nil
}

func (r *CharactersHandler) CreateMetricProperty(ctx context.Context, characterID primitive.ObjectID, metricID primitive.ObjectID, input model.MetricPropertyInput) (*repo.MetricProperty, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	metricProperty := repo.MetricProperty{
		ID: primitive.NewObjectID(),
	}

	if input.Name != nil {
		metricProperty.Name = *input.Name
	}
	if input.Type != nil {
		metricProperty.Type = (*input.Type)
	}
	if input.Value != nil {
		metricProperty.Value = *input.Value
	}
	if input.Unit != nil {
		metricProperty.Unit = *input.Unit
	}

	found := false
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricID {
			if len(character.CustomMetrics[i].Properties) >= int(character.CustomMetrics[i].LimitedPropertyNumber) {
				return nil, fmt.Errorf("metric properties creation limit reached")
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

	return &metricProperty, nil
}

func (r *CharactersHandler) UpdateMetricProperty(ctx context.Context, id primitive.ObjectID, characterID primitive.ObjectID, metricID primitive.ObjectID, input model.MetricPropertyInput) (*repo.MetricProperty, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	foundForMetric := false
	foundForProperty := false
	updatedProperty := repo.MetricProperty{}
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricID {
			for j, prop := range character.CustomMetrics[i].Properties {
				if prop.ID != id {
					continue
				}

				if input.Name != nil {
					prop.Name = *input.Name
				}
				if input.Type != nil {
					prop.Type = (*input.Type)
				}
				if input.Value != nil {
					prop.Value = *input.Value
				}
				if input.Unit != nil {
					prop.Unit = *input.Unit
				}

				character.CustomMetrics[i].Properties[j] = prop
				updatedProperty = prop
				foundForProperty = true
				break
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

	return &updatedProperty, nil
}

func (r *CharactersHandler) DeleteMetricProperty(ctx context.Context, id primitive.ObjectID, characterID primitive.ObjectID, metricID primitive.ObjectID) (*repo.MetricProperty, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	foundForMetric := false
	foundForProperty := false
	deletedMetricProperty := repo.MetricProperty{}
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricID {
			for j, prop := range character.CustomMetrics[i].Properties {
				if prop.ID == id {
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

	return &deletedMetricProperty, nil
}
