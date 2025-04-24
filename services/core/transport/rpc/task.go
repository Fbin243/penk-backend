package rpc

import (
	"context"

	"tenkhours/proto/pb/common"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (hdl *CoreHandler) UpsertTask(ctx context.Context, req *core.TaskInput) (*core.TaskMsg, error) {
	taskInput, err := Map[core.TaskInput, entity.TaskInput](req, UnixTimeConverter)
	if err != nil {
		return nil, err
	}

	task, err := hdl.taskBiz.Upsert(ctx, taskInput)
	if err != nil {
		return nil, err
	}

	return Map[entity.Task, core.TaskMsg](task, UnixTimeConverter)
}

func (hdl *CoreHandler) DeleteTask(ctx context.Context, req *common.IdReq) (*common.IdResp, error) {
	task, err := hdl.taskBiz.Delete(ctx, req.Id)

	return &common.IdResp{
		Id: task.ID,
	}, err
}

func (hdl *CoreHandler) UpsertTaskSession(ctx context.Context, req *core.TaskSessionInput) (*core.TaskSession, error) {
	taskSessionInput, err := Map[core.TaskSessionInput, entity.TaskSessionInput](req, UnixTimeConverter)
	if err != nil {
		return nil, err
	}

	taskSession, err := hdl.taskBiz.UpsertTaskSession(ctx, taskSessionInput)
	if err != nil {
		return nil, err
	}

	return Map[entity.TaskSession, core.TaskSession](taskSession, UnixTimeConverter)
}

func (hdl *CoreHandler) UpsertTaskSessions(ctx context.Context, req *core.TaskSessionInputs) (*core.TaskSessions, error) {
	entityTaskSessionInputs, err := MapSlice[core.TaskSessionInput, entity.TaskSessionInput](lo.FromSlicePtr(req.TaskSessionInputs), UnixTimeConverter)
	if err != nil {
		return nil, err
	}

	entityTaskSessions, err := hdl.taskBiz.UpsertTaskSessions(ctx, entityTaskSessionInputs)
	if err != nil {
		return nil, err
	}

	rpcTaskSessions, err := MapSlice[entity.TaskSession, core.TaskSession](entityTaskSessions, UnixTimeConverter)
	if err != nil {
		return nil, err
	}

	return &core.TaskSessions{
		TaskSessions: lo.ToSlicePtr(rpcTaskSessions),
	}, nil
}

func (hdl *CoreHandler) DeleteTaskSession(ctx context.Context, req *common.IdReq) (*common.IdResp, error) {
	taskSession, err := hdl.taskBiz.DeleteTaskSession(ctx, req.Id)

	return &common.IdResp{
		Id: taskSession.ID,
	}, err
}
