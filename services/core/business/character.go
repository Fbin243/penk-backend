package business

import (
	"context"
	"strconv"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"

	"github.com/samber/lo"
)

type CharacterBusiness struct {
	CharacterRepo ICharacterRepo
	ProfileRepo   IProfileRepo
	GoalRepo      IGoalRepo
}

type GoalTodo struct {
	updateMetrics    bool
	removeProperties bool
	checkFinish      bool
}

type (
	MetricMap   map[string]PropertyMap
	PropertyMap map[string]entity.MetricProperty
)

func NewCharacterBusiness(characterRepo ICharacterRepo, profileRepo IProfileRepo, goalRepo IGoalRepo) *CharacterBusiness {
	return &CharacterBusiness{
		CharacterRepo: characterRepo,
		ProfileRepo:   profileRepo,
		GoalRepo:      goalRepo,
	}
}

func (biz *CharacterBusiness) GetCharacterByID(ctx context.Context, id string) (*entity.Character, error) {
	character, err := biz.CharacterRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (biz *CharacterBusiness) GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	characters, err := biz.CharacterRepo.GetCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (biz *CharacterBusiness) UpsertCharacter(ctx context.Context, input entity.CharacterInput) (*entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	profile, err := biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	if input.ID == nil {
		// Insert new character
		charactersCount, err := biz.CharacterRepo.CountCharactersByProfileID(ctx, authSession.ProfileID)
		if err != nil {
			return nil, err
		}

		if charactersCount >= int64(profile.LimitedCharacterNumber) {
			return nil, errors.NewGQLError(errors.ErrCodeLimitCharacter, nil)
		}

		character := entity.Character{
			BaseEntity:          &base.BaseEntity{},
			Name:                input.Name,
			Gender:              input.Gender,
			ProfileID:           profile.ID,
			TotalFocusedTime:    0,
			CustomMetrics:       []entity.CustomMetric{},
			LimitedMetricNumber: utils.LimitedMetricNumber,
		}

		if input.Tags != nil {
			character.Tags = input.Tags
		}

		if input.Vision != nil {
			character.Vision = entity.Vision{
				Name: input.Vision.Name,
			}

			if input.Vision.Description != nil {
				character.Vision.Description = *input.Vision.Description
			}
		}

		if input.CustomMetrics != nil {
			err := biz.upsertMetricsInCharacter(ctx, &character, input.CustomMetrics)
			if err != nil {
				return nil, err
			}
		}

		createdCharacter, err := biz.CharacterRepo.InsertOne(ctx, &character)
		if err != nil {
			return nil, err
		}

		// Update the current character of the profile
		profile.CurrentCharacterID = character.ID
		_, err = biz.ProfileRepo.UpdateByID(ctx, profile.ID, profile)
		if err != nil {
			return nil, err
		}

		return createdCharacter, nil
	}

	// Update existing character
	character, err := biz.CharacterRepo.FindByID(ctx, *input.ID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != profile.ID {
		return nil, errors.PermissionDenied()
	}

	character.Name = input.Name

	if input.Tags != nil {
		character.Tags = input.Tags
	}

	if input.Vision != nil {
		character.Vision = entity.Vision{
			Name: input.Vision.Name,
		}

		if input.Vision.Description != nil {
			character.Vision.Description = *input.Vision.Description
		}
	}

	if input.CustomMetrics != nil {
		err := biz.upsertMetricsInCharacter(ctx, character, input.CustomMetrics)
		if err != nil {
			return nil, err
		}
	}

	return biz.CharacterRepo.UpdateByID(ctx, *input.ID, character)
}

func (biz *CharacterBusiness) upsertMetricsInCharacter(ctx context.Context, character *entity.Character, metricInputs []entity.CustomMetricInput) error {
	// Convert character metrics to map
	metricsMap := make(map[string]entity.CustomMetric)
	for _, metric := range character.CustomMetrics {
		metricsMap[metric.ID] = metric
	}

	metrics := make([]entity.CustomMetric, 0)

	// Check the limit metrics
	if len(metricInputs) > int(character.LimitedMetricNumber) {
		return errors.NewGQLError(errors.ErrCodeLimitMetric, nil)
	}

	// Get all unfinished, unexpired goals of a character
	goals, err := biz.GoalRepo.GetGoalsByCharacterID(ctx, character.ID, &entity.GoalStatusFilter{FinishStatus: lo.ToPtr(entity.GoalFinishStatusUnfinished), ExpireStatus: lo.ToPtr(entity.GoalExpireStatusUnexpired)})
	if err != nil {
		return err
	}

	goalIDs := lo.Map(goals, func(goal entity.Goal, _ int) string {
		return goal.ID
	})

	goalsTodo := &GoalTodo{
		updateMetrics:    false,
		removeProperties: false,
		checkFinish:      false,
	}

	for _, metricInput := range metricInputs {
		if metricInput.ID == nil {
			// Insert new metric
			metric := entity.CustomMetric{
				ID:                    mongodb.GenObjectID(),
				Name:                  metricInput.Name,
				Time:                  0,
				LimitedPropertyNumber: utils.LimitedPropertyNumber,
			}

			if metricInput.Description != nil {
				metric.Description = *metricInput.Description
			}

			if metricInput.Style != nil {
				metric.Style = entity.MetricStyle{
					Color: metricInput.Style.Color,
					Icon:  metricInput.Style.Icon,
				}
			}

			if metricInput.Properties != nil {
				err := biz.upsertPropertiesInMetric(ctx, &metric, metricInput.Properties, goals, goalsTodo)
				if err != nil {
					return err
				}
			}

			metrics = append(metrics, metric)
		} else {
			// Update existing metric
			updateMetricInGoal := false
			if _, ok := metricsMap[*metricInput.ID]; !ok {
				return errors.PermissionDenied()
			}

			existingMetric := metricsMap[*metricInput.ID]
			if existingMetric.Name != metricInput.Name {
				updateMetricInGoal = true
			}
			existingMetric.Name = metricInput.Name

			if metricInput.Description != nil {
				updateMetricInGoal = true
				existingMetric.Description = *metricInput.Description
			}

			if metricInput.Style != nil {
				updateMetricInGoal = true
				existingMetric.Style = entity.MetricStyle{
					Color: metricInput.Style.Color,
					Icon:  metricInput.Style.Icon,
				}
			}

			if metricInput.Properties != nil {
				err := biz.upsertPropertiesInMetric(ctx, &existingMetric, metricInput.Properties, goals, goalsTodo)
				if err != nil {
					return err
				}
			}

			if updateMetricInGoal {
				result, err := biz.GoalRepo.UpdateOneMetricInGoals(ctx, existingMetric, goalIDs)
				if err != nil {
					return err
				}
				if result.ModifiedCount > 0 {
					goalsTodo.updateMetrics = true
				}
			}

			metrics = append(metrics, existingMetric)
		}
	}

	if !goalsTodo.updateMetrics && !goalsTodo.removeProperties && goalsTodo.checkFinish {
		// Check if the goal is finished
		biz.checkGoalsFinished(ctx, goals, metrics, character.CustomMetrics)
	}

	character.CustomMetrics = metrics

	return nil
}

func (biz *CharacterBusiness) upsertPropertiesInMetric(ctx context.Context, metric *entity.CustomMetric, propertyInputs []entity.MetricPropertyInput, goals []entity.Goal, goalsTodo *GoalTodo) error {
	// Convert metric properties to map
	propertiesMap := make(map[string]entity.MetricProperty)
	for _, property := range metric.Properties {
		propertiesMap[property.ID] = property
	}

	// Check the limit properties
	if len(propertyInputs) > int(metric.LimitedPropertyNumber) {
		return errors.NewGQLError(errors.ErrCodeLimitProperty, nil)
	}

	properties := make([]entity.MetricProperty, 0)

	for _, propertyInput := range propertyInputs {
		if propertyInput.ID == nil {
			// Insert new property
			property := entity.MetricProperty{
				Name:  propertyInput.Name,
				Type:  propertyInput.Type,
				Value: propertyInput.Value,
				Unit:  propertyInput.Unit,
			}

			properties = append(properties, property)
		} else {
			removeProperty := false
			// Update existing property
			if _, ok := propertiesMap[*propertyInput.ID]; !ok {
				return errors.PermissionDenied()
			}

			property := propertiesMap[*propertyInput.ID]
			if property.Name != propertyInput.Name ||
				property.Type != propertyInput.Type ||
				property.Unit != propertyInput.Unit {
				removeProperty = true
			}
			property.Name = propertyInput.Name
			property.Type = propertyInput.Type
			property.Unit = propertyInput.Unit

			// Check if the value is changed, to check finish of the goal
			if property.Value != propertyInput.Value {
				goalsTodo.checkFinish = true
			}

			property.Value = propertyInput.Value

			properties = append(properties, property)

			if removeProperty {
				// Filter goals that include the property to be removed
				updateGoals := lo.Filter(goals, func(goal entity.Goal, _ int) bool {
					metrics := lo.Filter(goal.Target, func(_metric entity.CustomMetric, _ int) bool {
						return metric.ID == _metric.ID
					})

					if len(metrics) > 0 {
						props := lo.Filter(metrics[0].Properties, func(prop entity.MetricProperty, _ int) bool {
							return prop.ID == *propertyInput.ID
						})

						if len(props) > 0 {
							return true
						}
					}

					return false
				})

				updateGoalsIDs := lo.Map(updateGoals, func(goal entity.Goal, _ int) string {
					return goal.ID
				})

				result, err := biz.GoalRepo.RemoveOnePropertyFromGoals(ctx, metric.ID, *propertyInput.ID, updateGoalsIDs)
				if err != nil {
					return err
				}
				if result.ModifiedCount > 0 {
					goalsTodo.removeProperties = true
				}
			}
		}
	}

	metric.Properties = properties

	return nil
}

func (biz *CharacterBusiness) checkGoalsFinished(ctx context.Context, goals []entity.Goal, newMetrics, oldMetrics []entity.CustomMetric) error {
	// Convert new metrics to map
	newMetricsMap := MetricMap{}
	for _, metric := range newMetrics {
		newMetricsMap[metric.ID] = PropertyMap{}
		for _, property := range metric.Properties {
			newMetricsMap[metric.ID][property.ID] = property
		}
	}

	// Convert old metrics to map
	oldMetricsMap := MetricMap{}
	for _, metric := range oldMetrics {
		oldMetricsMap[metric.ID] = PropertyMap{}
		for _, property := range metric.Properties {
			oldMetricsMap[metric.ID][property.ID] = property
		}
	}

	finishedGoalIDs := make([]string, 0)
	// Check if the goal is finished
	for _, goal := range goals {
		// Compare goal target with new metrics and old metrics
		for _, metric := range goal.Target {
			newMetric, ok := newMetricsMap[metric.ID]
			if !ok {
				return errors.PermissionDenied()
			}

			oldMetric, ok := oldMetricsMap[metric.ID]
			if !ok {
				return errors.PermissionDenied()
			}

			// Compare properties
			for _, property := range metric.Properties {
				newProperty, ok := newMetric[property.ID]
				if !ok {
					return errors.PermissionDenied()
				}

				oldProperty, ok := oldMetric[property.ID]
				if !ok {
					return errors.PermissionDenied()
				}

				// Compare the property value based on the type
				switch property.Type {
				case entity.MetricPropertyTypeNumber:
					newValue, err := strconv.Atoi(newProperty.Value)
					if err != nil {
						return err
					}

					oldValue, err := strconv.Atoi(oldProperty.Value)
					if err != nil {
						return err
					}

					curValue, err := strconv.Atoi(property.Value)
					if err != nil {
						return err
					}

					if (newValue <= curValue && curValue <= oldValue) ||
						(oldValue <= curValue && curValue <= newValue) {
						goal.Status = entity.GoalFinishStatusFinished
						finishedGoalIDs = append(finishedGoalIDs, goal.ID)
					}
				case entity.MetricPropertyTypeString:
					if newProperty.Value == property.Value {
						goal.Status = entity.GoalFinishStatusFinished
						finishedGoalIDs = append(finishedGoalIDs, goal.ID)
					}
				}
			}
		}
	}

	// Update the status of the finished goals
	if len(finishedGoalIDs) > 0 {
		_, err := biz.GoalRepo.UpdateStatusOfGoals(ctx, finishedGoalIDs, entity.GoalFinishStatusFinished)
		if err != nil {
			return err
		}
	}

	return nil
}

func (biz *CharacterBusiness) DeleteCharacter(ctx context.Context, id string) (*entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	character, err := biz.CharacterRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != authSession.ProfileID {
		return nil, errors.PermissionDenied()
	}

	deletedCharacter, err := biz.CharacterRepo.DeleteCharacter(ctx, id)
	if err != nil {
		return nil, err
	}

	return deletedCharacter, nil
}

func (biz *CharacterBusiness) UpdateTimeInCharacter(ctx context.Context, characterID string, metricID *string, time int32) error {
	character, err := biz.CharacterRepo.FindByID(ctx, characterID)
	if err != nil {
		return err
	}

	character.TotalFocusedTime += time
	if metricID != nil {
		found := false
		for i, metric := range character.CustomMetrics {
			if metric.ID == *metricID {
				character.CustomMetrics[i].Time += time
				found = true
				break
			}
		}

		if !found {
			return errors.PermissionDenied()
		}
	}

	_, err = biz.CharacterRepo.UpdateByID(ctx, characterID, character)
	if err != nil {
		return err
	}

	return nil
}
