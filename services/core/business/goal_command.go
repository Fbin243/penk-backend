package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (biz *GoalBusiness) Upsert(ctx context.Context, input *entity.GoalInput) (*entity.Goal, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	goal := &entity.Goal{
		BaseEntity:  &base.BaseEntity{},
		CharacterID: authSession.CurrentCharacterID,
	}

	if input.ID == nil {
		count, err := biz.goalRepo.CountByCharacterID(ctx, authSession.CurrentCharacterID)
		if err != nil {
			return nil, err
		}

		if count >= utils.LimitedGoalNumber {
			return nil, errors.ErrLimitGoal
		}
	} else {
		err := biz.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
			{
				ID:   *input.ID,
				Type: entity.EntityTypeGoal,
			},
		})
		if err != nil {
			return nil, err
		}

		goal, err = biz.goalRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}
	}

	goal.Name = input.Name
	goal.StartTime = input.StartTime
	goal.EndTime = input.EndTime
	if input.Description != nil {
		goal.Description = *input.Description
	}

	// Get currrent metrics
	metrics, err := biz.metricRepo.Find(ctx, entity.MetricPipeline{
		Filter: &entity.MetricFilter{
			CharacterID: &authSession.CurrentCharacterID,
		},
	})
	if err != nil {
		return nil, err
	}

	metricMap := map[string]entity.Metric{}
	for _, metric := range metrics {
		metricMap[metric.ID] = metric
	}

	if input.Metrics != nil {
		err = biz.upsertGoalMetrics(ctx, upsertGoalMetricsParams{
			goal,
			input.Metrics,
			metricMap,
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

	if !goal.IsCompleted(metricMap) {
		goal.CompletedTime = nil
	} else if goal.CompletedTime == nil {
		goal.CompletedTime = lo.ToPtr(utils.Now())
	}

	if input.ID != nil {
		goal, err = biz.goalRepo.FindAndUpdateByID(ctx, *input.ID, goal)
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

type upsertGoalMetricsParams struct {
	goal         *entity.Goal
	metricInputs []entity.GoalMetricInput
	metricMap    map[string]entity.Metric
}

func (biz *GoalBusiness) upsertGoalMetrics(_ context.Context, params upsertGoalMetricsParams) error {
	metrics := []entity.GoalMetric{}
	for _, metricInput := range params.metricInputs {
		_, ok := params.metricMap[metricInput.ID]
		if !ok {
			return errors.NewGQLError(errors.ErrCodeNotFound, "Metric not found")
		}

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

	if params.goal != nil {
		params.goal.Metrics = metrics
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
				return errors.NewGQLError(errors.ErrCodeNotFound, "Checkbox not found")
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

func (biz *GoalBusiness) Delete(ctx context.Context, goalID string) (*entity.Goal, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	err = biz.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   goalID,
			Type: entity.EntityTypeGoal,
		},
	})
	if err != nil {
		return nil, err
	}

	goal, err := biz.goalRepo.FindAndDeleteByID(ctx, goalID)
	if err != nil {
		return nil, err
	}

	return goal, nil
}
