package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
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
	var character *entity.Character
	metricMap := map[string]entity.Metric{}

	if input.ID != nil {
		// Update existing goal
		goal, err = biz.GoalRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}

		if goal.CharacterID != input.CharacterID {
			return nil, errors.ErrPermissionDenied
		}

		// Just update if the goal is still unfinished and unexpired
		if goal.Status == entity.GoalFinishStatusFinished {
			return nil, errors.NewGQLError(errors.ErrCodeGoalAlreadyFinished, nil)
		}

		if goal.EndTime.Before(time.Now()) {
			return nil, errors.NewGQLError(errors.ErrCodeGoalAlreadyExpired, nil)
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
	if input.Metrics != nil {
		character, err = biz.CharacterRepo.FindByID(ctx, goal.CharacterID)
		if err != nil {
			return nil, err
		}

		for _, metric := range character.Metrics {
			metricMap[metric.ID] = metric
		}

		err := biz.upsertMetricsInGoal(ctx, upsertMetricInGoalInput{
			goal:         goal,
			metricInputs: input.Metrics,
			character:    character,
			metricMap:    metricMap,
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
	if input.Metrics != nil ||
		input.Checkboxes != nil {
		goal.UpdateStatus(metricMap)
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

type upsertMetricInGoalInput struct {
	goal         *entity.Goal
	metricInputs []entity.GoalMetricInput
	character    *entity.Character
	metricMap    map[string]entity.Metric
}

func (biz *GoalBusiness) upsertMetricsInGoal(_ context.Context, input upsertMetricInGoalInput) error {
	metrics := []entity.GoalTargetMetric{}
	for _, metric := range input.character.Metrics {
		input.metricMap[metric.ID] = metric
	}

	for _, metricInput := range input.metricInputs {
		if _, ok := input.metricMap[metricInput.ID]; !ok {
			return errors.NewGQLError(errors.ErrCodeBadRequest, "Metric not found")
		}

		metric := entity.GoalTargetMetric{
			ID:        input.metricMap[metricInput.ID].ID,
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

	input.goal.Target.Metrics = metrics

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
