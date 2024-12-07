package business

import (
	"context"
	"fmt"
	"strconv"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoalsBusiness struct {
	GoalsRepo      *repo.GoalsRepo
	CharactersRepo *repo.CharactersRepo
}

func NewGoalsBusiness(goalsRepo *repo.GoalsRepo, charactersRepo *repo.CharactersRepo) *GoalsBusiness {
	return &GoalsBusiness{goalsRepo, charactersRepo}
}

func (biz *GoalsBusiness) GetGoals(ctx context.Context, characterID primitive.ObjectID) ([]repo.Goal, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// Check if the character belongs to the user
	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorUnauthorized
	}

	return biz.GoalsRepo.GetGoalsByCharacterID(characterID)
}

func (biz *GoalsBusiness) UpsertGoal(ctx context.Context, characterID primitive.ObjectID, input model.GoalInput) (*repo.Goal, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// Check if the character belongs to the user
	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorUnauthorized
	}

	var newGoal *repo.Goal
	var existingGoal *repo.Goal
	updateFields := bson.M{
		"Name":      input.Name,
		"StartDate": input.StartDate,
		"EndDate":   input.EndDate,
	}

	if input.ID != nil {
		// Get the existing goal
		existingGoal, err = biz.GoalsRepo.FindByID(*input.ID)
		if err != nil {
			return nil, err
		}

		// Check permision
		if existingGoal.CharacterID != characterID {
			return nil, errors.ErrorPermissionDenied
		}

		// Just update if the goal is still active
		if existingGoal.Status != repo.GoalStatusActive {
			return nil, fmt.Errorf("goal is not active")
		}
	} else {
		// Create new goal
		newGoal = &repo.Goal{
			BaseModel:   &db.BaseModel{},
			CharacterID: characterID,
			Name:        input.Name,
			StartDate:   input.StartDate,
			EndDate:     input.EndDate,
			Status:      repo.GoalStatusActive,
		}
	}

	if input.Description != nil {
		newGoal.Description = *input.Description
		updateFields["Description"] = *input.Description
	}

	if len(input.Target) > 0 {
		targets := make([]repo.CustomMetric, 0)

		// Convert custom metrics to map for validation
		customMetricsMap := make(map[primitive.ObjectID]repo.CustomMetric)
		for _, metric := range character.CustomMetrics {
			customMetricsMap[metric.ID] = metric
		}

		// Add tracked metric to the goal
		for _, metric := range input.Target {
			// Validate metric in target
			if _, ok := customMetricsMap[metric.ID]; !ok {
				return nil, errors.ErrorPermissionDenied
			}

			currentMetric := customMetricsMap[metric.ID]

			trackedMetric := repo.CustomMetric{
				ID:          metric.ID,
				Name:        currentMetric.Name,
				Description: currentMetric.Description,
				Time:        currentMetric.Time,
				Style:       currentMetric.Style,
				Properties:  make([]repo.MetricProperty, 0),
			}

			// Convert properties to map for validation
			propertiesMap := make(map[primitive.ObjectID]repo.MetricProperty)
			for _, property := range currentMetric.Properties {
				propertiesMap[property.ID] = property
			}

			// Validate properties target metric
			for _, property := range metric.Properties {
				if _, ok := propertiesMap[property.ID]; !ok {
					return nil, errors.ErrorPermissionDenied
				}

				currentProperty := propertiesMap[property.ID]

				// Check if the property value is valid
				switch currentProperty.Type {
				case repo.MetricPropertyTypeNumber:
					_, err := strconv.Atoi(property.Value)
					if err != nil {
						return nil, err
					}
				}

				trackedMetric.Properties = append(trackedMetric.Properties, repo.MetricProperty{
					ID:    property.ID,
					Name:  currentProperty.Name,
					Type:  currentProperty.Type,
					Value: property.Value,
					Unit:  currentProperty.Unit,
				})
			}

			targets = append(targets, trackedMetric)
		}

		newGoal.Target = targets
		updateFields["Target"] = targets
	}

	if input.ID != nil {
		return biz.GoalsRepo.UpdateByID(*input.ID, updateFields)
	}

	return biz.GoalsRepo.InsertOne(newGoal)
}
