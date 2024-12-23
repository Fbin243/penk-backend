package business

import (
	"context"
	"strconv"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoalsBusiness struct {
	GoalsRepo      *repo.GoalsRepo
	CharactersRepo *repo.CharactersRepo
}

func NewGoalsBusiness(goalsRepo *repo.GoalsRepo, charactersRepo *repo.CharactersRepo) *GoalsBusiness {
	return &GoalsBusiness{goalsRepo, charactersRepo}
}

func (biz *GoalsBusiness) GetGoals(ctx context.Context, characterID primitive.ObjectID, status *repo.GoalStatusFilter) ([]repo.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	// Check if the character belongs to the user
	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != authSession.ProfileID {
		return nil, errors.PermissionDenied()
	}

	return biz.GoalsRepo.GetGoalsByCharacterID(characterID, status)
}

func (biz *GoalsBusiness) UpsertGoal(ctx context.Context, characterID primitive.ObjectID, input model.GoalInput) (*repo.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	// Check if the character belongs to the user
	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != authSession.ProfileID {
		return nil, errors.PermissionDenied()
	}

	var goal *repo.Goal

	if input.ID != nil {
		// Get the existing goal
		goal, err = biz.GoalsRepo.FindByID(*input.ID)
		if err != nil {
			return nil, err
		}

		// Check permision
		if goal.CharacterID != characterID {
			return nil, errors.PermissionDenied()
		}

		// Just update if the goal is still unfinished and unexpired
		if goal.Status == repo.GoalFinishStatusFinished {
			return nil, errors.NewError(errors.ErrCodeGoalAlreadyFinished, nil)
		}

		if goal.EndDate.Before(time.Now()) {
			return nil, errors.NewError(errors.ErrCodeGoalAlreadyExpired, nil)
		}
	} else {
		// Create new goal
		goal = &repo.Goal{
			BaseModel:   &db.BaseModel{},
			CharacterID: characterID,
			Name:        input.Name,
			StartDate:   input.StartDate,
			EndDate:     input.EndDate,
			Status:      repo.GoalFinishStatusUnfinished,
		}
	}

	if input.Description != nil {
		goal.Description = *input.Description
	}

	if input.Target != nil {
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
				return nil, errors.PermissionDenied()
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

			// Validate properties in a metric of the target
			for _, property := range metric.Properties {
				if _, ok := propertiesMap[property.ID]; !ok {
					return nil, errors.PermissionDenied()
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

		goal.Target = targets
	}

	if input.ID != nil {
		return biz.GoalsRepo.UpdateByID(*input.ID, goal)
	}

	return biz.GoalsRepo.InsertOne(goal)
}
