package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type TaskSession struct {
	*base.BaseEntity `                     bson:",inline"`
	TaskID           string     `json:"taskID"        bson:"task_id"`
	StartTime        time.Time  `json:"startTime"     bson:"start_time"`
	EndTime          time.Time  `json:"endTime"       bson:"end_time"`
	CompletedTime    *time.Time `json:"completedTime" bson:"completed_time"`
}
