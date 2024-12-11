package business

import (
	"context"
	"fmt"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharactersBusiness struct {
	CharactersRepo *repo.CharactersRepo
	ProfilesRepo   *repo.ProfilesRepo
	GoalsRepo      *repo.GoalsRepo
}

func NewCharactersBusiness(charactersRepo *repo.CharactersRepo, profilesRepo *repo.ProfilesRepo, goalsRepo *repo.GoalsRepo) *CharactersBusiness {
	return &CharactersBusiness{
		CharactersRepo: charactersRepo,
		ProfilesRepo:   profilesRepo,
		GoalsRepo:      goalsRepo,
	}
}

func (biz *CharactersBusiness) GetCharacterByID(ctx context.Context, id string) (*repo.Character, error) {
	characterOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	character, err := biz.CharactersRepo.FindByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to find character: %v", err)
	}

	return character, nil
}

func (biz *CharactersBusiness) GetCharactersByProfileID(ctx context.Context) ([]repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	characters, err := biz.CharactersRepo.GetCharactersByProfileID(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	return characters, nil
}

func (biz *CharactersBusiness) UpsertCharacter(ctx context.Context, input model.CharacterInput) (*repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	if input.ID == nil {
		charactersCount, err := biz.CharactersRepo.CountCharactersByProfileID(profile.ID)
		if err != nil {
			return nil, err
		}

		if charactersCount >= 2 {
			return nil, errors.ErrorCharacterLimitReached
		}

		character := repo.Character{
			BaseModel:           &db.BaseModel{},
			Name:                input.Name,
			Gender:              input.Gender,
			ProfileID:           profile.ID,
			TotalFocusedTime:    0,
			CustomMetrics:       []repo.CustomMetric{},
			LimitedMetricNumber: utils.LimitedMetricNumber,
		}

		if input.Tags != nil {
			character.Tags = input.Tags
		}

		if input.CustomMetrics != nil {
			err := biz.upsertMetricInCharacter(&character, input.CustomMetrics)
			if err != nil {
				return nil, err
			}
		}

		return biz.CharactersRepo.InsertOne(&character)
	}

	character, err := biz.CharactersRepo.FindByID(*input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find character: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
	}

	character.Name = input.Name

	if input.Tags != nil {
		character.Tags = input.Tags
	}

	if input.CustomMetrics != nil {
		err := biz.upsertMetricInCharacter(character, input.CustomMetrics)
		if err != nil {
			return nil, err
		}
	}

	return biz.CharactersRepo.UpdateByID(*input.ID, character)
}

func (biz *CharactersBusiness) upsertMetricInCharacter(character *repo.Character, metricInputs []model.CustomMetricInput) error {
	// Convert character metrics to map
	metricsMap := make(map[primitive.ObjectID]repo.CustomMetric)
	for _, metric := range character.CustomMetrics {
		metricsMap[metric.ID] = metric
	}

	metrics := make([]repo.CustomMetric, 0)

	for _, metricInput := range metricInputs {
		if metricInput.ID == nil {
			// Insert new metric
			metric := repo.CustomMetric{
				ID:                    primitive.NewObjectID(),
				Name:                  metricInput.Name,
				Time:                  0,
				LimitedPropertyNumber: utils.LimitedPropertyNumber,
			}

			if metricInput.Description != nil {
				metric.Description = *metricInput.Description
			}

			if metricInput.Style != nil {
				metric.Style = repo.MetricStyle{
					Color: metricInput.Style.Color,
					Icon:  metricInput.Style.Icon,
				}
			}

			if metricInput.Properties != nil {
				err := biz.upsertPropertyInMetric(&metric, metricInput.Properties)
				if err != nil {
					return err
				}
			}

			metrics = append(metrics, metric)
		} else {
			// Update existing metric
			if _, ok := metricsMap[*metricInput.ID]; !ok {
				return errors.ErrorPermissionDenied
			}

			metric := metricsMap[*metricInput.ID]
			metric.Name = metricInput.Name

			if metricInput.Description != nil {
				metric.Description = *metricInput.Description
			}

			if metricInput.Style != nil {
				metric.Style = repo.MetricStyle{
					Color: metricInput.Style.Color,
					Icon:  metricInput.Style.Icon,
				}
			}

			if metricInput.Properties != nil {
				err := biz.upsertPropertyInMetric(&metric, metricInput.Properties)
				if err != nil {
					return err
				}
			}

			metrics = append(metrics, metric)
		}
	}

	character.CustomMetrics = metrics

	return nil
}

func (biz *CharactersBusiness) upsertPropertyInMetric(metric *repo.CustomMetric, propertyInputs []model.MetricPropertyInput) error {
	// Convert metric properties to map
	propertiesMap := make(map[primitive.ObjectID]repo.MetricProperty)
	for _, property := range metric.Properties {
		propertiesMap[property.ID] = property
	}

	properties := make([]repo.MetricProperty, 0)

	for _, propertyInput := range propertyInputs {
		if propertyInput.ID == nil {
			// Insert new property
			property := repo.MetricProperty{
				ID:    primitive.NewObjectID(),
				Name:  propertyInput.Name,
				Type:  propertyInput.Type,
				Value: propertyInput.Value,
				Unit:  propertyInput.Unit,
			}

			properties = append(properties, property)
		} else {
			// Update existing property
			if _, ok := propertiesMap[*propertyInput.ID]; !ok {
				return errors.ErrorPermissionDenied
			}

			property := propertiesMap[*propertyInput.ID]
			property.Name = propertyInput.Name
			property.Type = propertyInput.Type
			property.Value = propertyInput.Value
			property.Unit = propertyInput.Unit

			properties = append(properties, property)
		}
	}

	metric.Properties = properties

	return nil
}

func (biz *CharactersBusiness) CreateCharacter(ctx context.Context, input model.CharacterInput) (*repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// TODO: Check if the user has already created 2 characters, maybe changed later
	characters, err := biz.CharactersRepo.GetCharactersByProfileID(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	if len(characters) >= 2 {
		return nil, fmt.Errorf("user have already created 2 characters")
	}

	character := repo.Character{
		BaseModel:           &db.BaseModel{},
		Name:                input.Name,
		Gender:              input.Gender,
		ProfileID:           profile.ID,
		TotalFocusedTime:    0,
		CustomMetrics:       []repo.CustomMetric{},
		LimitedMetricNumber: utils.LimitedMetricNumber,
	}

	if input.Tags != nil {
		character.Tags = input.Tags
	}

	// Create custom metrics for the character
	if input.CustomMetrics != nil {
		ctx = context.WithValue(ctx, FromCreateCharacter, true)
		for _, customMetric := range input.CustomMetrics {
			// Insert the character into context
			ctx := context.WithValue(ctx, CharacterKey, &character)
			biz.CreateCustomMetric(ctx, character.ID, customMetric)
		}
	}

	createdCharacter, err := biz.CharactersRepo.InsertOne(&character)
	if err != nil {
		return nil, fmt.Errorf("failed to create character: %v", err)
	}

	// TODO: Character has been created, so set the current character of the user to it
	profile.CurrentCharacterID = createdCharacter.ID
	_, err = biz.ProfilesRepo.UpdateProfile(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}

	return createdCharacter, nil
}

func (biz *CharactersBusiness) UpdateCharacter(ctx context.Context, id primitive.ObjectID, input model.CharacterInput) (*repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
	}

	character.Name = input.Name
	if input.Tags != nil {
		character.Tags = input.Tags
	}

	// Update custom metrics for the character
	if input.CustomMetrics != nil {
		ctx = context.WithValue(ctx, FromUpdateCharacter, true)
		// Insert the character into context
		ctx = context.WithValue(ctx, CharacterKey, character)

		metrics := make([]repo.CustomMetric, 0)
		for _, customMetric := range input.CustomMetrics {
			if customMetric.ID != nil {
				// Update custom metric
				metric, err := biz.UpdateCustomMetric(ctx, *customMetric.ID, character.ID, customMetric)
				if err != nil {
					return nil, fmt.Errorf("failed to update custom metric: %v", err)
				}

				metrics = append(metrics, *metric)
			} else {
				// Create custom metric
				metric, err := biz.CreateCustomMetric(ctx, character.ID, customMetric)
				if err != nil {
					return nil, fmt.Errorf("failed to create custom metric: %v", err)
				}

				metrics = append(metrics, *metric)
			}
		}

		character.CustomMetrics = metrics
	}

	updatedCharacter, err := biz.CharactersRepo.UpdateByID(id, character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %v", err)
	}

	return updatedCharacter, nil
}

func (biz *CharactersBusiness) DeleteCharacter(ctx context.Context, id primitive.ObjectID) (*repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
	}

	deletedCharacter, err := biz.CharactersRepo.DeleteCharacter(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete character: %v", err)
	}

	return deletedCharacter, nil
}

func (biz *CharactersBusiness) ResetCharacter(ctx context.Context, id primitive.ObjectID) (*repo.Character, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
	}

	character.Tags = []string{}
	character.TotalFocusedTime = 0
	character.CustomMetrics = []repo.CustomMetric{}

	resetCharacter, err := biz.CharactersRepo.UpdateByID(id, character)
	if err != nil {
		return nil, fmt.Errorf("failed to reset character: %v", err)
	}

	return resetCharacter, nil
}
