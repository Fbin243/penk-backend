package rpc

import (
	"context"

	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
)

func (hdl *CoreHandler) UpsertTask(ctx context.Context, req *core.TaskInput) (*core.TaskMsg, error) {
	taskInput, err := MapRPCInputToEntityInput[core.TaskInput, entity.TaskInput](req, nil)
	if err != nil {
		return nil, err
	}

	task, err := hdl.taskBiz.Upsert(ctx, taskInput)
	if err != nil {
		return nil, err
	}

	return MapEntityToRPC[entity.Task, core.TaskMsg](task, nil)
}

func (hdl *CoreHandler) UpsertTaskSession(ctx context.Context, req *core.TaskSessionInput) (*core.TaskSession, error) {
	taskSessionInput, err := MapRPCInputToEntityInput[core.TaskSessionInput, entity.TaskSessionInput](req, nil)
	if err != nil {
		return nil, err
	}

	taskSession, err := hdl.taskBiz.UpsertTaskSession(ctx, taskSessionInput)
	if err != nil {
		return nil, err
	}

	return MapEntityToRPC[entity.TaskSession, core.TaskSession](taskSession, nil)
}
