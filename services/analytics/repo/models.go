package repo

import (
	"time"

	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Metadata struct {
	CharacterID primitive.ObjectID `json:"characterID" bson:"character_id"`
	ProfileID   primitive.ObjectID `json:"profileID" bson:"profile_id"`
}

type Snapshot struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	Metadata    Metadata           `json:"metadata" bson:"metadata"`
	Character   repo.Character     `json:"character" bson:"character"`
	Description string             `json:"description,omitempty" bson:"description" validate:"omitempty,max=255"`
	Asset       interface{}        `json:"asset,omitempty" bson:"asset"`
}

type DFCapturedRecord struct {
	ID               string `dataframe:"id"`
	ProfileID        string `dataframe:"profile_id"`
	CharacterID      string `dataframe:"character_id"`
	Year             int    `dataframe:"year"`
	Month            int    `dataframe:"month"`
	Week             int    `dataframe:"week"`
	Day              int    `dataframe:"day"`
	TotalFocusedTime int    `dataframe:"total_focused_time"`
}

type DFCapturedRecordCustomMetric struct {
	ID       string `dataframe:"id"`
	MetricID string `dataframe:"metric_id"`
	Time     int    `dataframe:"time"`
}
