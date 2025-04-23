package rpc

import (
	"context"

	"tenkhours/proto/pb/common"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
)

func (hdl *CoreHandler) UpsertTask(ctx context.Context, req *core.TaskInput) (*core.TaskMsg, error) {
	taskInput, err := MapRPCInputToEntityInput[core.TaskInput, entity.TaskInput](req, UnixTimeConverter)
	if err != nil {
		return nil, err
	}

	task, err := hdl.taskBiz.Upsert(ctx, taskInput)
	if err != nil {
		return nil, err
	}

	return MapEntityToRPC[entity.Task, core.TaskMsg](task, UnixTimeConverter)
}

func (hdl *CoreHandler) DeleteTask(ctx context.Context, req *common.IdReq) (*common.IdResp, error) {
	task, err := hdl.taskBiz.Delete(ctx, req.Id)

	return &common.IdResp{
		Id: task.ID,
	}, err
}

func (hdl *CoreHandler) UpsertTaskSession(ctx context.Context, req *core.TaskSessionInput) (*core.TaskSession, error) {
	taskSessionInput, err := MapRPCInputToEntityInput[core.TaskSessionInput, entity.TaskSessionInput](req, UnixTimeConverter)
	if err != nil {
		return nil, err
	}

	taskSession, err := hdl.taskBiz.UpsertTaskSession(ctx, taskSessionInput)
	if err != nil {
		return nil, err
	}

	return MapEntityToRPC[entity.TaskSession, core.TaskSession](taskSession, UnixTimeConverter)
}

func (hdl *CoreHandler) UpsertTaskSessions(ctx context.Context, req *core.TaskSessionInputs) (*core.TaskSessions, error) {
	taskSessionInputs := []entity.TaskSessionInput{}
	for _, taskSessionInput := range req.TaskSessionInputs {
		taskSessionInputEntity, err := MapRPCInputToEntityInput[core.TaskSessionInput, entity.TaskSessionInput](taskSessionInput, UnixTimeConverter)
		if err != nil {
			return nil, err
		}

		taskSessionInputs = append(taskSessionInputs, *taskSessionInputEntity)
	}

	taskSessions, err := hdl.taskBiz.UpsertTaskSessions(ctx, taskSessionInputs)
	if err != nil {
		return nil, err
	}

	rpcTaskSessions := []*core.TaskSession{}
	for _, taskSession := range taskSessions {
		taskSessionMsg, err := MapEntityToRPC[entity.TaskSession, core.TaskSession](&taskSession, UnixTimeConverter)
		if err != nil {
			return nil, err
		}
		rpcTaskSessions = append(rpcTaskSessions, taskSessionMsg)
	}

	return &core.TaskSessions{
		TaskSessions: rpcTaskSessions,
	}, nil
}

func (hdl *CoreHandler) DeleteTaskSession(ctx context.Context, req *common.IdReq) (*common.IdResp, error) {
	taskSession, err := hdl.taskBiz.DeleteTaskSession(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &common.IdResp{
		Id: taskSession.ID,
	}, nil
}
