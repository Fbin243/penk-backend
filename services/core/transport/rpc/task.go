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

	task, err := hdl.taskBiz.UpsertTask(ctx, taskInput)
	if err != nil {
		return nil, err
	}

	return MapEntityToRPC[entity.Task, core.TaskMsg](task, nil)
}
