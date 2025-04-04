package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (biz *GoalBusiness) GetGoals(ctx context.Context, status *entity.GoalStatus) ([]entity.Goal, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	goals, err := biz.goalRepo.GetGoalsByCharacterID(ctx, authSession.CurrentCharacterID)
	if err != nil {
		return nil, err
	}

	if status != nil {
		goals = lo.Filter(goals, func(goal entity.Goal, _ int) bool {
			return goal.EvaluateStatus() == *status
		})
	}

	return goals, nil
}
