package business

import (
	"context"
	"strconv"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharactersBusiness struct {
	CharactersRepo *repo.CharactersRepo
	ProfilesRepo   *repo.ProfilesRepo
	GoalsRepo      *repo.GoalsRepo
}

type GoalsTodo struct {
	updateMetrics    bool
	removeProperties bool
	checkFinish      bool
}

type MetricMap map[primitive.ObjectID]PropertyMap
type PropertyMap map[primitive.ObjectID]repo.MetricProperty

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
		return nil, err
	}

	return character, nil
}

func (biz *CharactersBusiness) GetCharactersByProfileID(ctx context.Context) ([]repo.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	characters, err := biz.CharactersRepo.GetCharactersByProfileID(authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (biz *CharactersBusiness) UpsertCharacter(ctx context.Context, input model.CharacterInput) (*repo.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	profile, err := biz.ProfilesRepo.FindByID(authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	if input.ID == nil {
		// Insert new character
		charactersCount, err := biz.CharactersRepo.CountCharactersByProfileID(authSession.ProfileID)
		if err != nil {
			return nil, err
		}

		if charactersCount >= int64(profile.LimitedCharacterNumber) {
			return nil, errors.NewError(errors.ErrCodeLimitCharacter, nil)
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

		if input.Vision != nil {
			character.Vision = repo.Vision{
				Name: input.Vision.Name,
			}

			if input.Vision.Description != nil {
				character.Vision.Description = *input.Vision.Description
			}
		}

		if input.CustomMetrics != nil {
			err := biz.upsertMetricsInCharacter(&character, input.CustomMetrics)
			if err != nil {
				return nil, err
			}
		}

		createdCharacter, err := biz.CharactersRepo.InsertOne(&character)
		if err != nil {
			return nil, err
		}

		// TODO: Character has been created, so set the current character of the user to it
		profile.CurrentCharacterID = character.ID
		_, err = biz.ProfilesRepo.UpdateByID(profile.ID, profile)
		if err != nil {
			return nil, err
		}

		return createdCharacter, nil
	}

	// Update existing character
	character, err := biz.CharactersRepo.FindByID(*input.ID)
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
		character.Vision = repo.Vision{
			Name: input.Vision.Name,
		}

		if input.Vision.Description != nil {
			character.Vision.Description = *input.Vision.Description
		}
	}

	if input.CustomMetrics != nil {
		err := biz.upsertMetricsInCharacter(character, input.CustomMetrics)
		if err != nil {
			return nil, err
		}
	}

	return biz.CharactersRepo.UpdateByID(*input.ID, character)
}

func (biz *CharactersBusiness) upsertMetricsInCharacter(character *repo.Character, metricInputs []model.CustomMetricInput) error {
	// Convert character metrics to map
	metricsMap := make(map[primitive.ObjectID]repo.CustomMetric)
	for _, metric := range character.CustomMetrics {
		metricsMap[metric.ID] = metric
	}

	metrics := make([]repo.CustomMetric, 0)

	// Check the limit metrics
	if len(metricInputs) > int(character.LimitedMetricNumber) {
		return errors.NewError(errors.ErrCodeLimitMetric, nil)
	}

	// Get all unfinished, unexpired goals of a character
	goals, err := biz.GoalsRepo.GetGoalsByCharacterID(character.ID, &repo.GoalStatusFilter{FinishStatus: lo.ToPtr(repo.GoalFinishStatusUnfinished), ExpireStatus: lo.ToPtr(repo.GoalExpireStatusUnexpired)})
	if err != nil {
		return err
	}

	goalIDs := lo.Map(goals, func(goal repo.Goal, _ int) primitive.ObjectID {
		return goal.ID
	})

	goalsTodo := &GoalsTodo{
		updateMetrics:    false,
		removeProperties: false,
		checkFinish:      false,
	}

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
				err := biz.upsertPropertiesInMetric(&metric, metricInput.Properties, goals, goalsTodo)
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
				existingMetric.Style = repo.MetricStyle{
					Color: metricInput.Style.Color,
					Icon:  metricInput.Style.Icon,
				}
			}

			if metricInput.Properties != nil {
				err := biz.upsertPropertiesInMetric(&existingMetric, metricInput.Properties, goals, goalsTodo)
				if err != nil {
					return err
				}
			}

			if updateMetricInGoal {
				result, err := biz.GoalsRepo.UpdateOneMetricInGoals(existingMetric, goalIDs)
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
		biz.checkGoalsFinished(goals, metrics, character.CustomMetrics)
	}

	character.CustomMetrics = metrics

	return nil
}

func (biz *CharactersBusiness) upsertPropertiesInMetric(metric *repo.CustomMetric, propertyInputs []model.MetricPropertyInput, goals []repo.Goal, goalsTodo *GoalsTodo) error {
	// Convert metric properties to map
	propertiesMap := make(map[primitive.ObjectID]repo.MetricProperty)
	for _, property := range metric.Properties {
		propertiesMap[property.ID] = property
	}

	// Check the limit properties
	if len(propertyInputs) > int(metric.LimitedPropertyNumber) {
		return errors.NewError(errors.ErrCodeLimitProperty, nil)
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
				updateGoals := lo.Filter(goals, func(goal repo.Goal, _ int) bool {
					metrics := lo.Filter(goal.Target, func(_metric repo.CustomMetric, _ int) bool {
						return metric.ID == _metric.ID
					})

					if len(metrics) > 0 {
						props := lo.Filter(metrics[0].Properties, func(prop repo.MetricProperty, _ int) bool {
							return prop.ID == *propertyInput.ID
						})

						if len(props) > 0 {
							return true
						}
					}

					return false
				})

				updateGoalsIDs := lo.Map(updateGoals, func(goal repo.Goal, _ int) primitive.ObjectID {
					return goal.ID
				})

				result, err := biz.GoalsRepo.RemoveOnePropertyFromGoals(metric.ID, *propertyInput.ID, updateGoalsIDs)
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

func (biz *CharactersBusiness) checkGoalsFinished(goals []repo.Goal, newMetrics []repo.CustomMetric, oldMetrics []repo.CustomMetric) error {
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

	finishedGoalIDs := make([]primitive.ObjectID, 0)
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
				case repo.MetricPropertyTypeNumber:
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
						goal.Status = repo.GoalFinishStatusFinished
						finishedGoalIDs = append(finishedGoalIDs, goal.ID)
					}
				case repo.MetricPropertyTypeString:
					if newProperty.Value == property.Value {
						goal.Status = repo.GoalFinishStatusFinished
						finishedGoalIDs = append(finishedGoalIDs, goal.ID)
					}
				}
			}
		}
	}

	// Update the status of the finished goals
	if len(finishedGoalIDs) > 0 {
		_, err := biz.GoalsRepo.UpdateStatusOfGoals(finishedGoalIDs, repo.GoalFinishStatusFinished)
		if err != nil {
			return err
		}
	}

	return nil
}

func (biz *CharactersBusiness) DeleteCharacter(ctx context.Context, id primitive.ObjectID) (*repo.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	character, err := biz.CharactersRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != authSession.ProfileID {
		return nil, errors.PermissionDenied()
	}

	deletedCharacter, err := biz.CharactersRepo.DeleteCharacter(id)
	if err != nil {
		return nil, err
	}

	return deletedCharacter, nil
}

func (biz *CharactersBusiness) UpdateTimeInCharacter(ctx context.Context, characterID primitive.ObjectID, metricID primitive.ObjectID, time int32) error {
	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return err
	}

	character.TotalFocusedTime += time
	if !metricID.IsZero() {
		found := false
		for i, metric := range character.CustomMetrics {
			if metric.ID == metricID {
				character.CustomMetrics[i].Time += time
				found = true
				break
			}
		}

		if !found {
			return errors.PermissionDenied()
		}
	}

	_, err = biz.CharactersRepo.UpdateByID(characterID, character)
	if err != nil {
		return err
	}

	return nil
}
