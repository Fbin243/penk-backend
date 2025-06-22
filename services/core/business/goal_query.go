package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (biz *GoalBusiness) Get(ctx context.Context, filter *entity.GoalFilter, orderBy *entity.GoalOrderBy, limit, offset *int) ([]entity.Goal, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = &entity.GoalFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	goals, err := biz.goalRepo.Find(ctx, entity.GoalPipeline{
		Filter:  filter,
		OrderBy: orderBy,
		Pagination: &types.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, err
	}

	if filter.Status != nil {
		goals = lo.Filter(goals, func(goal entity.Goal, _ int) bool {
			return goal.EvaluateStatus() == *filter.Status
		})
	}

	return goals, nil
}

func (biz *GoalBusiness) Count(ctx context.Context, filter *entity.GoalFilter) (int, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return 0, err
	}

	if filter == nil {
		filter = &entity.GoalFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	if filter.Status != nil {
		goals, err := biz.goalRepo.Find(ctx, entity.GoalPipeline{
			Filter: filter,
		})
		if err != nil {
			return 0, err
		}

		return len(lo.Filter(goals, func(goal entity.Goal, _ int) bool {
			return goal.EvaluateStatus() == *filter.Status
		})), nil
	}

	return biz.goalRepo.CountByFilter(ctx, filter)
}
