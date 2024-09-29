package core

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/core/validations"
	"tenkhours/pkg/db/coredb"
	"tenkhours/services/core_v2/graph/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *CharactersHandler) CreateCustomMetric(ctx context.Context, characterID string, input model.CustomMetricInput) (*coredb.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object id: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
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

	err = validations.ValidateCustomMetric(customMetric)
	if err != nil {
		return nil, err
	}

	createdCustomMetric, err := r.CharactersRepo.CreateCustomMetric(character.ID, &customMetric)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom metric: %v", err)
	}

	return createdCustomMetric, nil
}

func (r *CharactersHandler) UpdateCustomMetric(ctx context.Context, id string, characterID string, input model.CustomMetricInput) (*coredb.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	metricOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	found := false
	updatedMetric := coredb.CustomMetric{}
	for i, cm := range character.CustomMetrics {
		if cm.ID != metricOID {
			continue
		}
		if input.Name != nil {
			cm.Name = *input.Name
		}
		if input.Description != nil {
			cm.Description = *input.Description
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
			var properties []coredb.MetricProperty
			for j, prop := range input.Properties {
				var metricProperty coredb.MetricProperty
				if len(cm.Properties) > j {
					metricProperty.ID = cm.Properties[j].ID
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

				err = validations.ValidateMetricProperty(metricProperty)
				if err != nil {
					return nil, err
				}

				if len(properties) >= int(character.CustomMetrics[i].LimitedPropertyNumber) {
					return nil, fmt.Errorf("metric properties creation limit reached")
				}

				properties = append(properties, metricProperty)
			}
			cm.Properties = properties
		}

		err = validations.ValidateCustomMetric(cm)
		if err != nil {
			return nil, err
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

func (r *CharactersHandler) DeleteCustomMetric(ctx context.Context, id string, characterID string) (*coredb.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character id: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	metricOID, err := primitive.ObjectIDFromHex(id)
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

	return deletedCustomMetric, nil
}

func (r *CharactersHandler) ResetCustomMetric(ctx context.Context, id string, characterID string) (*coredb.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	metricOID, err := primitive.ObjectIDFromHex(id)
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

	return &resetMetric, nil
}

func (r *CharactersHandler) CreateMetricProperty(ctx context.Context, characterID string, metricID string, input model.MetricPropertyInput) (*coredb.MetricProperty, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	metricOID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	metricProperty := coredb.MetricProperty{
		ID: primitive.NewObjectID(),
	}

	if input.Name != nil {
		metricProperty.Name = *input.Name
	}
	if input.Type != nil {
		metricProperty.Type = *input.Type
	}
	if input.Value != nil {
		metricProperty.Value = *input.Value
	}
	if input.Unit != nil {
		metricProperty.Unit = *input.Unit
	}

	found := false
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricOID {
			if len(character.CustomMetrics[i].Properties) >= int(character.CustomMetrics[i].LimitedPropertyNumber) {
				return nil, fmt.Errorf("metric properties creation limit reached")
			}

			err = validations.ValidateMetricProperty(metricProperty)
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

	return &metricProperty, nil
}

func (r *CharactersHandler) UpdateMetricProperty(ctx context.Context, id string, characterID string, metricID string, input model.MetricPropertyInput) (*coredb.MetricProperty, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	metricOID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	metricPropOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid metric property ID: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	foundForMetric := false
	foundForProperty := false
	updatedProperty := coredb.MetricProperty{}
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricOID {
			for j, prop := range character.CustomMetrics[i].Properties {
				if prop.ID != metricPropOID {
					continue
				}

				if input.Name != nil {
					prop.Name = *input.Name
				}
				if input.Type != nil {
					prop.Type = *input.Type
				}
				if input.Value != nil {
					prop.Value = *input.Value
				}
				if input.Unit != nil {
					prop.Unit = *input.Unit
				}

				err := validations.ValidateMetricProperty(prop)
				if err != nil {
					return nil, err
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

func (r *CharactersHandler) DeleteMetricProperty(ctx context.Context, id string, characterID string, metricID string) (*coredb.MetricProperty, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, fmt.Errorf("invalid character ID: %v", err)
	}

	metricOID, err := primitive.ObjectIDFromHex(metricID)
	if err != nil {
		return nil, fmt.Errorf("invalid metric ID: %v", err)
	}

	metricPropOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid metric property ID: %v", err)
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
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

	return &deletedMetricProperty, nil
}
