package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
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
		return nil, errors.ErrUnauthorized
	}

	err := biz.CharacterRepo.ValidateCharacter(ctx, authSession.ProfileID, characterID)
	if err != nil {
		return nil, err
	}

	goals, err := biz.GoalRepo.GetGoalsByCharacterID(ctx, characterID, status)
	if err != nil {
		return nil, err
	}

	characterMap, err := biz.getCharacterMap(ctx, characterID)
	if err != nil {
		return nil, err
	}

	lo.ForEach(goals, func(goal entity.Goal, _ int) {
		if goal.Status == entity.GoalFinishStatusUnfinished {
			biz.populateGoal(ctx, &goal, characterMap)
		}
	})

	return goals, nil
}

func (biz *GoalBusiness) UpsertGoal(ctx context.Context, input entity.GoalInput) (*entity.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := biz.CharacterRepo.ValidateCharacter(ctx, authSession.ProfileID, input.CharacterID)
	if err != nil {
		return nil, err
	}

	var goal *entity.Goal
	characterMap, err := biz.getCharacterMap(ctx, input.CharacterID)
	if err != nil {
		return nil, err
	}

	if input.ID != nil {
		// Update existing goal
		goal, err = biz.GoalRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}

		if goal.CharacterID != input.CharacterID {
			return nil, errors.ErrPermissionDenied
		}

		if goal.Status == entity.GoalFinishStatusFinished {
			return nil, errors.NewGQLError(errors.ErrCodeGoalAlreadyFinished, nil)
		}
	} else {
		// Create new goal
		goal = &entity.Goal{
			BaseEntity:  &base.BaseEntity{},
			CharacterID: input.CharacterID,
			Status:      entity.GoalFinishStatusUnfinished,
		}
	}

	goal.Name = input.Name
	goal.StartTime = input.StartTime
	goal.EndTime = input.EndTime
	if input.Description != nil {
		goal.Description = *input.Description
	}

	if input.Categories != nil {
		err := biz.upsertGoalCategories(ctx, goal, input.Categories)
		if err != nil {
			return nil, err
		}
	}

	if input.Metrics != nil {
		err := biz.upsertGoalMetrics(ctx, upsertGoalMetricsParam{
			goal:         goal,
			metricInputs: input.Metrics,
		})
		if err != nil {
			return nil, err
		}
	}

	if input.Checkboxes != nil {
		err := biz.upsertCheckboxesInGoal(ctx, goal, input.Checkboxes)
		if err != nil {
			return nil, err
		}
	}

	if input.Categories != nil ||
		input.Metrics != nil ||
		input.Checkboxes != nil {
		goal.UpdateStatus(characterMap.MetricMap)
	}

	if goal.Status == entity.GoalFinishStatusFinished {
		biz.populateGoal(ctx, goal, characterMap)
	}

	if input.ID != nil {
		goal, err = biz.GoalRepo.UpdateByID(ctx, *input.ID, goal)
		if err != nil {
			return nil, err
		}
	} else {
		goal, err = biz.GoalRepo.InsertOne(ctx, goal)
		if err != nil {
			return nil, err
		}
	}

	return goal, nil
}

func (biz *GoalBusiness) upsertGoalCategories(ctx context.Context, goal *entity.Goal, categoryInputs []entity.GoalCategoryInput) error {
	categories := []entity.GoalCategory{}
	for _, categoryInput := range categoryInputs {
		category := entity.GoalCategory{
			Category: &entity.Category{
				ID: categoryInput.ID,
			},
		}

		if categoryInput.Metrics != nil {
			err := biz.upsertGoalMetrics(ctx,
				upsertGoalMetricsParam{
					category:     &category,
					metricInputs: categoryInput.Metrics,
				})
			if err != nil {
				return err
			}
		}

		categories = append(categories, category)
	}

	goal.Target.Categories = categories

	return nil
}

type upsertGoalMetricsParam struct {
	goal         *entity.Goal
	category     *entity.GoalCategory
	metricInputs []entity.GoalMetricInput
}

func (biz *GoalBusiness) upsertGoalMetrics(_ context.Context, param upsertGoalMetricsParam) error {
	metrics := []entity.GoalMetric{}
	for _, metricInput := range param.metricInputs {
		metric := entity.GoalMetric{
			Metric: &entity.Metric{
				ID: metricInput.ID,
			},
			Condition: metricInput.Condition,
		}

		switch metricInput.Condition {
		case entity.MetricConditionInRange:
			metric.RangeValue = &entity.Range{
				Min: metricInput.RangeValue.Min,
				Max: metricInput.RangeValue.Max,
			}
		default:
			metric.TargetValue = metricInput.TargetValue
		}

		metrics = append(metrics, metric)
	}

	if param.goal != nil {
		param.goal.Target.Metrics = metrics
	}
	if param.category != nil {
		param.category.Metrics = metrics
	}

	return nil
}

func (biz *GoalBusiness) upsertCheckboxesInGoal(_ context.Context, goal *entity.Goal, checkboxInputs []entity.CheckboxInput) error {
	checkboxes := []entity.Checkbox{}
	checkboxMap := map[string]entity.Checkbox{}
	for _, checkbox := range goal.Target.Checkboxes {
		checkboxMap[checkbox.ID] = checkbox
	}

	for _, checkboxInput := range checkboxInputs {
		checkbox := entity.Checkbox{}
		if checkboxInput.ID != nil {
			if _, ok := checkboxMap[*checkboxInput.ID]; !ok {
				return errors.NewGQLError(errors.ErrCodeBadRequest, "Category not found")
			}

			checkbox = checkboxMap[*checkboxInput.ID]
			checkbox.Name = checkboxInput.Name
			checkbox.Value = checkboxInput.Value
		} else {
			checkbox = entity.Checkbox{
				ID:    mongodb.GenObjectID(),
				Name:  checkboxInput.Name,
				Value: checkboxInput.Value,
			}
		}

		checkboxes = append(checkboxes, checkbox)
	}

	goal.Target.Checkboxes = checkboxes

	return nil
}

func (biz *GoalBusiness) DeleteGoal(ctx context.Context, goalID string) (*entity.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := biz.GoalRepo.ValidateGoal(ctx, authSession.ProfileID, goalID)
	if err != nil {
		return nil, err
	}

	goal, err := biz.GoalRepo.DeleteByID(ctx, goalID)
	if err != nil {
		return nil, err
	}

	return goal, nil
}

func (biz *GoalBusiness) getCharacterMap(ctx context.Context, characterID string) (entity.CharacterMap, error) {
	characterMap := entity.CharacterMap{
		CategoryMap: map[string]entity.Category{},
		MetricMap:   map[string]entity.Metric{},
	}
	character, err := biz.CharacterRepo.FindByID(ctx, characterID)
	if err != nil {
		return characterMap, err
	}

	for _, category := range character.Categories {
		characterMap.CategoryMap[category.ID] = category
		for _, metric := range category.Metrics {
			characterMap.MetricMap[metric.ID] = metric
		}
	}
	for _, metric := range character.Metrics {
		characterMap.MetricMap[metric.ID] = metric
	}

	return characterMap, nil
}

func (biz *GoalBusiness) populateGoal(_ context.Context, goal *entity.Goal, characterMap entity.CharacterMap) {
	goal.Target.Categories = lo.Map(goal.Target.Categories,
		func(category entity.GoalCategory, _ int) entity.GoalCategory {
			category.Category = lo.ToPtr(characterMap.CategoryMap[category.ID])

			category.Metrics = lo.Map(category.Metrics,
				func(metric entity.GoalMetric, _ int) entity.GoalMetric {
					metric.Metric = lo.ToPtr(characterMap.MetricMap[metric.ID])
					return metric
				})

			return category
		})

	goal.Target.Metrics = lo.Map(goal.Target.Metrics,
		func(metric entity.GoalMetric, _ int) entity.GoalMetric {
			metric.Metric = lo.ToPtr(characterMap.MetricMap[metric.ID])
			return metric
		})

	return
}
