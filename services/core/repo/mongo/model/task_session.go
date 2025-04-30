package mongomodel

import (
	"time"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskSession struct {
	*mongodb.BaseEntity `                     bson:",inline"`
	TaskOID             primitive.ObjectID `json:"taskID"        bson:"task_id"`
	StartTime           time.Time          `json:"startTime"     bson:"start_time"`
	EndTime             time.Time          `json:"endTime"       bson:"end_time"`
	CompletedTime       *time.Time         `json:"completedTime" bson:"completed_time"`
}

func (t *TaskSession) TaskID(id string) {
	t.TaskOID = mongodb.ToObjectID(id)
}
