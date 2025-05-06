package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DevicesToken struct {
	*mongodb.BaseEntity `                           bson:",inline"`
	ProfileID           string  `json:"profileID,omitempty" bson:"profile_id"`
	Tokens              []Token `json:"tokens,omitempty"    bson:"tokens"`
}

type Token struct {
	DeviceID string `json:"deviceID,omitempty" bson:"device_id"`
	Token    string `json:"token,omitempty"    bson:"token"`
	Platform string `json:"platform,omitempty" bson:"platform"`
	CreateAt string `json:"createAt,omitempty" bson:"create_at"`
}

type Reminder struct {
	*mongodb.BaseEntity `                            bson:",inline"`
	ProfileID           string             `json:"profile_id,omitempty" bson:"profile_id"`
	Type                ReminderType       `json:"type,omitempty"       bson:"type"`
	Title               string             `json:"title,omitempty"      bson:"title"`
	RemindTime          time.Time          `json:"remind_time,omitempty" bson:"remind_time"`
	Recurrence          string             `json:"recurrence,omitempty"  bson:"recurrence"`
	LinkedItemID        primitive.ObjectID `json:"linked_item_id,omitempty" bson:"linked_item_id"`
}

type ReminderType string

const (
	ReminderTypeTask  ReminderType = "TASK"
	ReminderTypeEvent ReminderType = "EVENT"
	ReminderTypeHabit ReminderType = "HABIT"
)
