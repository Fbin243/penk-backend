package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (biz *CharactersBusiness) CreateCustomMetric(ctx context.Context, characterID primitive.ObjectID, input model.CustomMetricInput) (*repo.CustomMetric, error) {
	var character *repo.Character
	var err error

	// Check if the request is from CreateCharacter or UpdateCharacter
	fromCreateCharacter, ok := ctx.Value(FromCreateCharacter).(bool)
	if !ok {
		fromCreateCharacter = false
	}

	fromUpdateCharacter, ok := ctx.Value(FromUpdateCharacter).(bool)
	if !ok {
		fromUpdateCharacter = false
	}

	if fromCreateCharacter || fromUpdateCharacter {
		character = ctx.Value(CharacterKey).(*repo.Character)
	} else {
		profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
		if !ok {
			return nil, errors.ErrorUnauthorized
		}

		character, err = biz.CharactersRepo.FindByID(characterID)
		if err != nil {
			return nil, fmt.Errorf("failed to get character: %v", err)
		}

		if character.ProfileID != profile.ID {
			return nil, errors.ErrorPermissionDenied
		}

		if len(character.CustomMetrics) >= int(character.LimitedMetricNumber) {
			return nil, fmt.Errorf("custom metric creation limit reached")
		}
	}

	customMetric := repo.CustomMetric{
		ID:                    primitive.NewObjectID(),
		Description:           "",
		Time:                  0,
		Style:                 repo.MetricStyle{},
		Properties:            []repo.MetricProperty{},
		LimitedPropertyNumber: utils.LimitedPropertyNumber,
	}

	customMetric.Name = input.Name
	if input.Description != nil {
		customMetric.Description = *input.Description
	}

	if input.Style != nil {
		customMetric.Style = repo.MetricStyle{
			Color: input.Style.Color,
			Icon:  input.Style.Icon,
		}
	}

	if input.Properties != nil {
		var properties []repo.MetricProperty
		for _, prop := range input.Properties {
			var metricProperty repo.MetricProperty
			metricProperty.ID = primitive.NewObjectID()

			metricProperty.Name = prop.Name
			metricProperty.Type = prop.Type
			metricProperty.Value = prop.Value
			metricProperty.Unit = prop.Unit

			if len(properties) >= int(customMetric.LimitedPropertyNumber) {
				return nil, errors.ErrorMetricLimitReached
			}

			properties = append(properties, metricProperty)
		}

		customMetric.Properties = properties
	}

	if fromCreateCharacter || fromUpdateCharacter {
		character.CustomMetrics = append(character.CustomMetrics, customMetric)
		return &customMetric, nil
	}

	createdCustomMetric, err := biz.CharactersRepo.CreateCustomMetric(character.ID, &customMetric)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom metric: %v", err)
	}

	return createdCustomMetric, nil
}

func (biz *CharactersBusiness) UpdateCustomMetric(ctx context.Context, metricID primitive.ObjectID, characterID primitive.ObjectID, input model.CustomMetricInput) (*repo.CustomMetric, error) {
	var character *repo.Character
	var err error

	// Check if the request is from UpdateCharacter
	fromUpdateCharacter, ok := ctx.Value(FromCreateCharacter).(bool)
	if !ok {
		fromUpdateCharacter = false
	}

	if fromUpdateCharacter {
		character = ctx.Value(CharacterKey).(*repo.Character)
	} else {
		profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
		if !ok {
			return nil, errors.ErrorUnauthorized
		}

		character, err = biz.CharactersRepo.FindByID(characterID)
		if err != nil {
			return nil, fmt.Errorf("failed to get character: %v", err)
		}

		if character.ProfileID != profile.ID {
			return nil, errors.ErrorPermissionDenied
		}
	}

	updateMetrics := lo.Filter(character.CustomMetrics, func(cm repo.CustomMetric, _ int) bool {
		return cm.ID == metricID
	})

	if len(updateMetrics) == 0 {
		return nil, errors.ErrorMetricNotFound
	}

	updateMetric := updateMetrics[0]
	updateMetric.Name = input.Name
	if input.Description != nil {
		updateMetric.Description = *input.Description
	}

	if input.Style != nil {
		updateMetric.Style = repo.MetricStyle{
			Color: input.Style.Color,
			Icon:  input.Style.Icon,
		}
	}

	if input.Properties != nil {
		var properties []repo.MetricProperty
		for _, prop := range input.Properties {
			var metricProperty repo.MetricProperty
			if prop.ID != nil {
				// Validate the property ID
				currentProps := lo.Filter(updateMetric.Properties, func(p repo.MetricProperty, _ int) bool {
					return p.ID == *prop.ID
				})

				if len(currentProps) == 0 {
					return nil, errors.ErrorPropertyNotFound
				}

				metricProperty = currentProps[0]
			} else {
				metricProperty.ID = primitive.NewObjectID()
			}

			metricProperty.Name = prop.Name
			metricProperty.Type = prop.Type
			metricProperty.Value = prop.Value
			metricProperty.Unit = prop.Unit

			if len(properties) >= int(updateMetric.LimitedPropertyNumber) {
				return nil, errors.ErrorMetricLimitReached
			}

			properties = append(properties, metricProperty)
		}

		updateMetric.Properties = properties
	}

	if fromUpdateCharacter {
		return &updateMetric, nil
	}

	_, err = biz.CharactersRepo.UpdateByID(characterID, character)
	if err != nil {
		return nil, fmt.Errorf("failed to update custom metric: %v", err)
	}

	return &updateMetric, nil
}

func (biz *CharactersBusiness) DeleteCustomMetric(ctx context.Context, metricID primitive.ObjectID, characterID primitive.ObjectID) (*repo.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
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

	deletedCustomMetric, err := biz.CharactersRepo.DeleteCustomMetric(characterID, metricID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete metric: %v", err)
	}

	return deletedCustomMetric, nil
}

func (biz *CharactersBusiness) ResetCustomMetric(ctx context.Context, metricID primitive.ObjectID, characterID primitive.ObjectID) (*repo.CustomMetric, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
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

	_, err = biz.CharactersRepo.UpdateByID(characterID, character)
	if err != nil {
		return nil, fmt.Errorf("failed to reset custom metric: %v", err)
	}

	return &resetMetric, nil
}

func (biz *CharactersBusiness) CreateMetricProperty(ctx context.Context, characterID primitive.ObjectID, metricID primitive.ObjectID, input model.MetricPropertyInput) (*repo.MetricProperty, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
	}

	metricProperty := repo.MetricProperty{
		ID: primitive.NewObjectID(),
	}

	metricProperty.Name = input.Name
	metricProperty.Type = input.Type
	metricProperty.Value = input.Value
	metricProperty.Unit = input.Unit

	found := false
	for i, cm := range character.CustomMetrics {
		if cm.ID == metricID {
			if len(character.CustomMetrics[i].Properties) >= int(character.CustomMetrics[i].LimitedPropertyNumber) {
				return nil, errors.ErrorMetricLimitReached
			}

			character.CustomMetrics[i].Properties = append(character.CustomMetrics[i].Properties, metricProperty)
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("custom metric does not belong to the character")
	}

	_, err = biz.CharactersRepo.UpdateByID(characterID, character)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric property: %v", err)
	}

	return &metricProperty, nil
}

func (biz *CharactersBusiness) UpdateMetricProperty(ctx context.Context, id primitive.ObjectID, characterID primitive.ObjectID, metricID primitive.ObjectID, input model.MetricPropertyInput) (*repo.MetricProperty, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
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

				prop.Name = input.Name
				prop.Type = input.Type
				prop.Value = input.Value
				prop.Unit = input.Unit

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

	_, err = biz.CharactersRepo.UpdateByID(characterID, character)
	if err != nil {
		return nil, fmt.Errorf("failed to update metric property: %v", err)
	}

	return &updatedProperty, nil
}

func (biz *CharactersBusiness) DeleteMetricProperty(ctx context.Context, id primitive.ObjectID, characterID primitive.ObjectID, metricID primitive.ObjectID) (*repo.MetricProperty, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
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

	_, err = biz.CharactersRepo.UpdateByID(characterID, character)
	if err != nil {
		return nil, fmt.Errorf("failed to delete metric property: %v", err)
	}

	return &deletedMetricProperty, nil
}
