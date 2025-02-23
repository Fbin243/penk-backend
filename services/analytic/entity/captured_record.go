package entity

import "time"

type CapturedRecord struct {
	ID               string                       `json:"id"                      bson:"_id,omitempty"`
	Timestamp        time.Time                    `json:"timestamp"               bson:"timestamp,omitempty"`
	TotalFocusedTime int32                        `json:"totalFocusedTime"        bson:"total_focused_time,omitempty"`
	Categories       []CapturedRecordCategory     `json:"categories,omitempty"    bson:"categories,omitempty"`
	TimeTrackings    []CapturedRecordTimeTracking `json:"timeTrackings,omitempty" bson:"time_trackings,omitempty"`
	Metadata         CapturedRecordMetadata       `json:"metadata"                bson:"metadata,omitempty"`
}

type CapturedRecordCategory struct {
	ID   string `json:"id"   bson:"_id,omitempty"`
	Time int32  `json:"time" bson:"time,omitempty"`
}

type CapturedRecordMetadata struct {
	CharacterID string `json:"characterID" bson:"character_id,omitempty"`
	ProfileID   string `json:"profileID"   bson:"profile_id,omitempty"`
}

type CapturedRecordTimeTracking struct {
	CategoryID *string   `json:"categoryID,omitempty" bson:"category_id,omitempty"`
	Time       int32     `json:"time"                 bson:"time,omitempty"`
	StartTime  time.Time `json:"startTime"            bson:"start_time,omitempty"`
	EndTime    time.Time `json:"endTime"              bson:"end_time,omitempty"`
}

type GetCapturedRecordFilter struct {
	ProfileID   string
	CharacterID *string
	StartTime   time.Time
	EndTime     time.Time
}
