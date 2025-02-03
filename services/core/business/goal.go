package business

import (
	"context"
	"strconv"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"
)

type GoalBusiness struct {
	GoalRepo      IGoalRepo
	CharacterRepo ICharacterRepo
}

func NewGoalBusiness(goalRepo IGoalRepo, characterRepo ICharacterRepo) *GoalBusiness {
	return &GoalBusiness{goalRepo, characterRepo}
}

func (biz *GoalBusiness) GetGoals(ctx context.Context, characterID string, status *entity.GoalStatusFilter) ([]entity.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	// Check if the character belongs to the user
	character, err := biz.CharacterRepo.FindByID(ctx, characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != authSession.ProfileID {
		return nil, errors.PermissionDenied()
	}

	return biz.GoalRepo.GetGoalsByCharacterID(ctx, characterID, status)
}

func (biz *GoalBusiness) UpsertGoal(ctx context.Context, characterID string, input entity.GoalInput) (*entity.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	// Check if the character belongs to the user
	character, err := biz.CharacterRepo.FindByID(ctx, characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != authSession.ProfileID {
		return nil, errors.PermissionDenied()
	}

	var goal *entity.Goal

	if input.ID != nil {
		// Get the existing goal
		goal, err = biz.GoalRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}

		// Check permission
		if goal.CharacterID != characterID {
			return nil, errors.PermissionDenied()
		}

		// Just update if the goal is still unfinished and unexpired
		if goal.Status == entity.GoalFinishStatusFinished {
			return nil, errors.NewGQLError(errors.ErrCodeGoalAlreadyFinished, nil)
		}

		if goal.EndDate.Before(time.Now()) {
			return nil, errors.NewGQLError(errors.ErrCodeGoalAlreadyExpired, nil)
		}
	} else {
		// Create new goal
		goal = &entity.Goal{
			BaseEntity:  &base.BaseEntity{},
			CharacterID: characterID,
			Name:        input.Name,
			StartDate:   input.StartDate,
			EndDate:     input.EndDate,
			Status:      entity.GoalFinishStatusUnfinished,
		}
	}

	if input.Description != nil {
		goal.Description = *input.Description
	}

	if input.Target != nil {
		targets := make([]entity.CustomMetric, 0)

		// Convert custom metrics to map for validation
		customMetricsMap := make(map[string]entity.CustomMetric)
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

			trackedMetric := entity.CustomMetric{
				ID:          metric.ID,
				Name:        currentMetric.Name,
				Description: currentMetric.Description,
				Time:        currentMetric.Time,
				Style:       currentMetric.Style,
				Properties:  make([]entity.MetricProperty, 0),
			}

			// Convert properties to map for validation
			propertiesMap := make(map[string]entity.MetricProperty)
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
				case entity.MetricPropertyTypeNumber:
					_, err := strconv.Atoi(property.Value)
					if err != nil {
						return nil, err
					}
				}

				trackedMetric.Properties = append(trackedMetric.Properties, entity.MetricProperty{
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
		return biz.GoalRepo.UpdateByID(ctx, *input.ID, goal)
	}

	return biz.GoalRepo.InsertOne(ctx, goal)
}
