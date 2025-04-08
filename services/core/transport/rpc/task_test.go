package rpc_test

import (
	"log"
	"testing"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
	"tenkhours/services/core/transport/rpc"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// Entity
var taskCheckbox = entity.Checkbox{
	ID:    mongodb.GenObjectID(),
	Name:  "Subtask name",
	Value: true,
}

var task = entity.Task{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	},
	CharacterID:   mongodb.GenObjectID(),
	CategoryID:    lo.ToPtr(mongodb.GenObjectID()),
	Name:          "Task name",
	Priority:      1,
	CompletedTime: lo.ToPtr(utils.Now()),
	Subtasks: []entity.Checkbox{
		taskCheckbox, taskCheckbox,
	},
	Description: lo.ToPtr("Task description"),
	Deadline:    lo.ToPtr(utils.Now()),
}

func TestMapTask(t *testing.T) {
	rpcTask, err := rpc.MapEntityToRPC[entity.Task, core.TaskMsg](&task, rpc.UnixTimeConverter)
	assert.NoError(t, err)

	log.Printf("Entity Task: %+v", utils.PrettyJSON(task))
	log.Printf("RPC Task: %+v", utils.PrettyJSON(rpcTask))

	assert.Equal(t, task.ID, rpcTask.Id)
	assert.Equal(t, task.CharacterID, rpcTask.CharacterId)
	assert.Equal(t, task.CategoryID, rpcTask.CategoryId)
	assert.Equal(t, task.Name, rpcTask.Name)
	assert.Equal(t, task.Priority, int(rpcTask.Priority))
	assert.Equal(t, task.CompletedTime.Unix(), *rpcTask.CompletedTime)
	assert.Equal(t, task.CreatedAt.Unix(), rpcTask.CreatedAt)
	assert.Equal(t, task.UpdatedAt.Unix(), rpcTask.UpdatedAt)
	assert.Equal(t, task.Description, rpcTask.Description)
	assert.Equal(t, task.Deadline.Unix(), *rpcTask.Deadline)

	for i, expectedCheckbox := range task.Subtasks {
		assertCheckbox(t, expectedCheckbox, rpcTask.Subtasks[i])
	}
}

// RPC Input
var rpcTaskCheckboxInput = &core.CheckboxInput{
	Id:    lo.ToPtr(mongodb.GenObjectID()),
	Name:  "Subtask name",
	Value: true,
}

var rpcTaskInput = &core.TaskInput{
	Id:            lo.ToPtr(mongodb.GenObjectID()),
	CategoryId:    lo.ToPtr(mongodb.GenObjectID()),
	Name:          "Task name",
	Priority:      1,
	CompletedTime: lo.ToPtr(utils.Now().Unix()),
	Subtasks: []*core.CheckboxInput{
		rpcTaskCheckboxInput, rpcTaskCheckboxInput,
	},
	Description: lo.ToPtr("Task description"),
	Deadline:    lo.ToPtr(utils.Now().Unix()),
}

func TestMapTaskInput(t *testing.T) {
	taskInput, err := rpc.MapRPCInputToEntityInput[core.TaskInput, entity.TaskInput](rpcTaskInput, rpc.UnixTimeConverter)
	assert.NoError(t, err)

	log.Printf("RPC Task Input: %+v", utils.PrettyJSON(rpcTaskInput))
	log.Printf("Entity Task Input: %+v", utils.PrettyJSON(taskInput))

	assert.Equal(t, rpcTaskInput.Id, taskInput.ID)
	assert.Equal(t, rpcTaskInput.CategoryId, taskInput.CategoryID)
	assert.Equal(t, rpcTaskInput.Name, taskInput.Name)
	assert.Equal(t, rpcTaskInput.Priority, int32(taskInput.Priority))
	assert.Equal(t, *rpcTaskInput.CompletedTime, taskInput.CompletedTime.Unix())
	assert.Equal(t, rpcTaskInput.Description, taskInput.Description)
	assert.Equal(t, *rpcTaskInput.Deadline, taskInput.Deadline.Unix())

	for i, expectedCheckbox := range rpcTaskInput.Subtasks {
		assertCheckboxInput(t, expectedCheckbox, taskInput.Subtasks[i])
	}
}
