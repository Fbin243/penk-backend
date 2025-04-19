package rpc

import (
	"context"

	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
)

func (hdl *CoreHandler) UpsertGoal(ctx context.Context, req *core.GoalInput) (*core.Goal, error) {
	goalInput, err := MapRPCInputToEntityInput[core.GoalInput, entity.GoalInput](req, append(UnixTimeConverter, MetricConditionConverter...))
	if err != nil {
		return nil, err
	}

	goal, err := hdl.goalBiz.Upsert(ctx, goalInput)
	if err != nil {
		return nil, err
	}

	return MapEntityToRPC[entity.Goal, core.Goal](goal, append(UnixTimeConverter, MetricConditionConverter...))
}
