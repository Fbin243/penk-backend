package entity

import "time"

type CapturedRecord struct {
	ID               string                       `json:"id"                      bson:"_id,omitempty"`
	Timestamp        time.Time                    `json:"timestamp"               bson:"timestamp,omitempty"`
	TotalFocusedTime int32                        `json:"totalFocusedTime"        bson:"total_focused_time,omitempty"`
	CustomMetrics    []CapturedRecordCustomMetric `json:"categories,omitempty"    bson:"custom_metrics,omitempty"`
	TimeTrackings    []CapturedRecordTimeTracking `json:"timeTrackings,omitempty" bson:"time_trackings,omitempty"`
	Metadata         CapturedRecordMetadata       `json:"metadata"                bson:"metadata,omitempty"`
}

type CapturedRecordCustomMetric struct {
	ID   string `json:"id"   bson:"_id,omitempty"`
	Time int32  `json:"time" bson:"time,omitempty"`
}

type CapturedRecordMetadata struct {
	CharacterID string `json:"characterID" bson:"character_id,omitempty"`
	ProfileID   string `json:"profileID"   bson:"profile_id,omitempty"`
}

type CapturedRecordTimeTracking struct {
	CustomMetricID string    `json:"customMetricID" bson:"custom_metric_id,omitempty"`
	Time           int32     `json:"time"           bson:"time,omitempty"`
	StartTime      time.Time `json:"startTime"      bson:"start_time,omitempty"`
	EndTime        time.Time `json:"endTime"        bson:"end_time,omitempty"`
}
