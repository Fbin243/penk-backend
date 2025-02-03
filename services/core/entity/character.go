package entity

import "tenkhours/pkg/db/base"

type Character struct {
	*base.BaseEntity    `                                     bson:",inline"`
	ProfileID           string         `json:"profileID,omitempty"           bson:"profile_id"`
	Name                string         `json:"name,omitempty"                bson:"name"`
	Gender              bool           `json:"gender,omitempty"              bson:"gender"`
	Tags                []string       `json:"tags,omitempty"                bson:"tags"`
	TotalFocusedTime    int32          `json:"totalFocusedTime,omitempty"    bson:"total_focused_time"`
	CustomMetrics       []CustomMetric `json:"customMetrics,omitempty"       bson:"custom_metrics"`
	LimitedMetricNumber int32          `json:"limitedMetricNumber,omitempty" bson:"limited_metric_number"`
	Vision              Vision         `json:"vision,omitempty"              bson:"vision"`
}

type Vision struct {
	Name        string `json:"name,omitempty"        bson:"name"`
	Description string `json:"description,omitempty" bson:"description"`
}

type CustomMetric struct {
	ID                    string           `json:"id,omitempty"                    bson:"_id"`
	Name                  string           `json:"name,omitempty"                  bson:"name"`
	Description           string           `json:"description,omitempty"           bson:"description"`
	Time                  int32            `json:"time,omitempty"                  bson:"time"`
	Style                 MetricStyle      `json:"style,omitempty"                 bson:"style"`
	Properties            []MetricProperty `json:"properties,omitempty"            bson:"properties"`
	LimitedPropertyNumber int32            `json:"limitedPropertyNumber,omitempty" bson:"limited_property_number"`
}

type MetricStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty"  bson:"icon"`
}

type MetricProperty struct {
	ID    string             `json:"id,omitempty"    bson:"_id"`
	Name  string             `json:"name,omitempty"  bson:"name"`
	Type  MetricPropertyType `json:"type,omitempty"  bson:"type"`
	Value string             `json:"value,omitempty" bson:"value"`
	Unit  string             `json:"unit,omitempty"  bson:"unit"`
}

type MetricPropertyType string

const (
	MetricPropertyTypeNumber MetricPropertyType = "Number"
	MetricPropertyTypeString MetricPropertyType = "String"
)
