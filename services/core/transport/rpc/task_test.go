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
	rpcTask, err := rpc.Map[entity.Task, core.TaskMsg](&task, rpc.UnixTimeConverter)
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
	taskInput, err := rpc.Map[core.TaskInput, entity.TaskInput](rpcTaskInput, rpc.UnixTimeConverter)
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

// -------------------- Task session --------------------
// Entity
var taskSession = entity.TaskSession{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	},
	TaskID:        mongodb.GenObjectID(),
	StartTime:     utils.Now(),
	EndTime:       utils.Now(),
	CompletedTime: lo.ToPtr(utils.Now()),
}

func TestMapTaskSession(t *testing.T) {
	rpcTaskSession, err := rpc.Map[entity.TaskSession, core.TaskSession](&taskSession, rpc.UnixTimeConverter)
	assert.NoError(t, err)

	log.Printf("Entity Task Session: %+v", utils.PrettyJSON(taskSession))
	log.Printf("RPC Task Session: %+v", utils.PrettyJSON(rpcTaskSession))

	assert.Equal(t, taskSession.ID, rpcTaskSession.Id)
	assert.Equal(t, taskSession.TaskID, rpcTaskSession.TaskId)
	assert.Equal(t, taskSession.StartTime.Unix(), rpcTaskSession.StartTime)
	assert.Equal(t, taskSession.EndTime.Unix(), rpcTaskSession.EndTime)
	assert.Equal(t, taskSession.CompletedTime.Unix(), *rpcTaskSession.CompletedTime)
}

// RPC Input
var rpcTaskSessionInput = &core.TaskSessionInput{
	Id:            lo.ToPtr(mongodb.GenObjectID()),
	TaskId:        mongodb.GenObjectID(),
	StartTime:     utils.Now().Unix(),
	EndTime:       utils.Now().Unix(),
	CompletedTime: lo.ToPtr(utils.Now().Unix()),
}

func TestMapTaskSessionInput(t *testing.T) {
	taskSessionInput, err := rpc.Map[core.TaskSessionInput, entity.TaskSessionInput](rpcTaskSessionInput, rpc.UnixTimeConverter)
	assert.NoError(t, err)

	log.Printf("RPC Task Session Input: %+v", utils.PrettyJSON(rpcTaskSessionInput))
	log.Printf("Entity Task Session Input: %+v", utils.PrettyJSON(taskSessionInput))

	assert.Equal(t, rpcTaskSessionInput.Id, taskSessionInput.ID)
	assert.Equal(t, rpcTaskSessionInput.TaskId, taskSessionInput.TaskID)
	assert.Equal(t, rpcTaskSessionInput.StartTime, taskSessionInput.StartTime.Unix())
	assert.Equal(t, rpcTaskSessionInput.EndTime, taskSessionInput.EndTime.Unix())
	assert.Equal(t, *rpcTaskSessionInput.CompletedTime, taskSessionInput.CompletedTime.Unix())
}
