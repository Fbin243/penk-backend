package rpc

import (
	"context"

	"tenkhours/proto/pb/common"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
)

func (hdl *CoreHandler) UpsertHabit(ctx context.Context, req *core.HabitInput) (*core.Habit, error) {
	converters := append(UnixTimeConverter, HabitConverter...)
	habitInput, err := Map[core.HabitInput, entity.HabitInput](req, converters)
	if err != nil {
		return nil, err
	}

	habit, err := hdl.habitBiz.Upsert(ctx, habitInput)
	if err != nil {
		return nil, err
	}

	return Map[entity.Habit, core.Habit](habit, converters)
}

func (hdl *CoreHandler) DeleteHabit(ctx context.Context, req *common.IdReq) (*common.IdResp, error) {
	category, err := hdl.habitBiz.Delete(ctx, req.Id)

	return &common.IdResp{
		Id: category.ID,
	}, err
}

func (hdl *CoreHandler) UpsertHabitLog(ctx context.Context, req *core.HabitLogInput) (*core.HabitLog, error) {
	habitLogInput, err := Map[core.HabitLogInput, entity.HabitLogInput](req, nil)
	if err != nil {
		return nil, err
	}

	habitLog, err := hdl.habitBiz.UpsertHabitLog(ctx, habitLogInput)
	if err != nil {
		return nil, err
	}

	return Map[entity.HabitLog, core.HabitLog](habitLog, nil)
}
