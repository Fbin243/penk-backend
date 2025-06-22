package rpc

import (
	"context"

	"tenkhours/proto/pb/common"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
)

func (hdl *CoreHandler) UpsertGoal(ctx context.Context, req *core.GoalInput) (*core.Goal, error) {
	goalInput, err := Map[core.GoalInput, entity.GoalInput](req, append(UnixTimeConverter, MetricConditionConverter...))
	if err != nil {
		return nil, err
	}

	goal, err := hdl.goalBiz.Upsert(ctx, goalInput)
	if err != nil {
		return nil, err
	}

	return Map[entity.Goal, core.Goal](goal, append(UnixTimeConverter, MetricConditionConverter...))
}

func (hdl *CoreHandler) DeleteGoal(ctx context.Context, req *common.IdReq) (*common.IdResp, error) {
	goal, err := hdl.goalBiz.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &common.IdResp{
		Id: goal.ID,
	}, nil
}
