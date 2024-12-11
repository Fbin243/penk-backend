package repo

import (
	"time"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	*db.BaseModel      `bson:",inline"`
	Name               string             `json:"name,omitempty" bson:"name"`
	Email              string             `json:"email,omitempty" bson:"email"`
	FirebaseUID        string             `json:"firebaseUID,omitempty" bson:"firebase_uid"`
	ImageURL           string             `json:"imageURL,omitempty" bson:"image_url"`
	CurrentCharacterID primitive.ObjectID `json:"currentCharacterID,omitempty" bson:"current_character_id"`
	AvailableSnapshots int32              `json:"availableSnapshots,omitempty" bson:"available_snapshots"`
	AutoSnapshot       bool               `json:"autoSnapshot,omitempty" bson:"auto_snapshot"`
}

type Character struct {
	*db.BaseModel       `bson:",inline"`
	ProfileID           primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id"`
	Name                string             `json:"name,omitempty" bson:"name"`
	Gender              bool               `json:"gender,omitempty" bson:"gender"`
	Tags                []string           `json:"tags,omitempty" bson:"tags"`
	TotalFocusedTime    int32              `json:"totalFocusedTime,omitempty" bson:"total_focused_time"`
	CustomMetrics       []CustomMetric     `json:"customMetrics,omitempty" bson:"custom_metrics"`
	LimitedMetricNumber int32              `json:"limitedMetricNumber,omitempty" bson:"limited_metric_number,omitempty"`
}

type CustomMetric struct {
	ID                    primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name                  string             `json:"name,omitempty" bson:"name"`
	Description           string             `json:"description,omitempty" bson:"description"`
	Time                  int32              `json:"time,omitempty" bson:"time,omitempty"`
	Style                 MetricStyle        `json:"style,omitempty" bson:"style"`
	Properties            []MetricProperty   `json:"properties,omitempty" bson:"properties"`
	LimitedPropertyNumber int32              `json:"limitedPropertyNumber,omitempty" bson:"limited_property_number,omitempty"`
}

type MetricStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty" bson:"icon"`
}

type MetricProperty struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name  string             `json:"name,omitempty" bson:"name"`
	Type  MetricPropertyType `json:"type,omitempty" bson:"type"`
	Value string             `json:"value,omitempty" bson:"value"`
	Unit  string             `json:"unit,omitempty" bson:"unit"`
}

type MetricPropertyType string

const (
	MetricPropertyTypeNumber MetricPropertyType = "Number"
	MetricPropertyTypeString MetricPropertyType = "String"
)

type Goal struct {
	*db.BaseModel `bson:",inline"`
	CharacterID   primitive.ObjectID `json:"characterID" bson:"character_id"`
	Name          string             `json:"name" bson:"name"`
	Description   string             `json:"description" bson:"description"`
	StartDate     time.Time          `json:"startDate" bson:"start_date"`
	EndDate       time.Time          `json:"endDate" bson:"end_date"`
	Status        GoalFinishStatus   `json:"status" bson:"status"`
	Target        []CustomMetric     `json:"target" bson:"target"`
}

type GoalFinishStatus string
type GoalExpireStatus string

const (
	GoalFinishStatusFinished   GoalFinishStatus = "Unfinnished"
	GoalFinishStatusUnfinished GoalFinishStatus = "Finished"
	GoalExpireStatusExpired    GoalExpireStatus = "Expired"
	GoalExpireStatusUnexpired  GoalExpireStatus = "Unexpired"
)

type GoalStatusFilter struct {
	FinishStatus *GoalFinishStatus
	ExpireStatus *GoalExpireStatus
}
