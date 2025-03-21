package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"
)

type GoalBusiness struct {
	goalRepo      IGoalRepo
	characterRepo ICharacterRepo
	categoryRepo  ICategoryRepo
	metricRepo    IMetricRepo
}

func NewGoalBusiness(goalRepo IGoalRepo, characterRepo ICharacterRepo, categoryRepo ICategoryRepo, metricRepo IMetricRepo) *GoalBusiness {
	return &GoalBusiness{goalRepo, characterRepo, categoryRepo, metricRepo}
}

func (biz *GoalBusiness) GetGoals(ctx context.Context, characterID string, status *entity.GoalStatus) ([]entity.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := biz.characterRepo.Exist(ctx, authSession.ProfileID, characterID)
	if err != nil {
		return nil, err
	}

	goals, err := biz.goalRepo.GetGoalsByCharacterID(ctx, characterID, status)
	if err != nil {
		return nil, err
	}

	return goals, nil
}

func (biz *GoalBusiness) UpsertGoal(ctx context.Context, input entity.GoalInput) (*entity.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := biz.characterRepo.Exist(ctx, authSession.ProfileID, input.CharacterID)
	if err != nil {
		return nil, err
	}

	var goal *entity.Goal
	if input.ID != nil {
		// Update existing goal
		goal, err = biz.goalRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}

		if goal.CharacterID != input.CharacterID {
			return nil, errors.ErrPermissionDenied
		}
	} else {
		// Create new goal
		goal = &entity.Goal{
			BaseEntity:  &base.BaseEntity{},
			CharacterID: input.CharacterID,
			Status:      entity.GoalStatusPlanned,
		}
	}

	goal.Name = input.Name
	goal.StartTime = input.StartTime
	goal.EndTime = input.EndTime
	if input.Description != nil {
		goal.Description = *input.Description
	}

	if input.Status != nil {
		goal.Status = *input.Status
	}

	if input.Metrics != nil {
		err := biz.upsertGoalMetrics(ctx, goal, input.Metrics)
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

	if input.Metrics != nil ||
		input.Checkboxes != nil {
		// Get currrent metrics
		metrics, err := biz.metricRepo.FindByCharacterID(ctx, goal.CharacterID)
		if err != nil {
			return nil, err
		}

		metricMap := map[string]entity.Metric{}
		for _, metric := range metrics {
			metricMap[metric.ID] = metric
		}
		goal.UpdateStatus(metricMap)
	}

	if input.ID != nil {
		goal, err = biz.goalRepo.UpdateByID(ctx, *input.ID, goal)
		if err != nil {
			return nil, err
		}
	} else {
		goal, err = biz.goalRepo.InsertOne(ctx, goal)
		if err != nil {
			return nil, err
		}
	}

	return goal, nil
}

func (biz *GoalBusiness) upsertGoalMetrics(_ context.Context, goal *entity.Goal, metricInputs []entity.GoalMetricInput) error {
	metrics := []entity.GoalMetric{}
	for _, metricInput := range metricInputs {
		metric := entity.GoalMetric{
			ID:        metricInput.ID,
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

	if goal != nil {
		goal.Metrics = metrics
	}

	return nil
}

func (biz *GoalBusiness) upsertCheckboxesInGoal(_ context.Context, goal *entity.Goal, checkboxInputs []entity.CheckboxInput) error {
	if len(checkboxInputs) > utils.LimitedCheckboxNumber {
		return errors.ErrLimitCheckbox
	}

	checkboxes := []entity.Checkbox{}
	checkboxMap := map[string]entity.Checkbox{}
	for _, checkbox := range goal.Checkboxes {
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

	goal.Checkboxes = checkboxes

	return nil
}

func (biz *GoalBusiness) DeleteGoal(ctx context.Context, goalID string) (*entity.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := biz.goalRepo.ValidateGoal(ctx, authSession.ProfileID, goalID)
	if err != nil {
		return nil, err
	}

	goal, err := biz.goalRepo.DeleteByID(ctx, goalID)
	if err != nil {
		return nil, err
	}

	return goal, nil
}
