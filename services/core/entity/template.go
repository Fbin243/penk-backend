package entity

import "tenkhours/pkg/db/base"

type Template struct {
	*base.BaseEntity `                             bson:",inline"`
	Name             string             `json:"name,omitempty"        bson:"name"`
	Description      string             `json:"description,omitempty" bson:"description"`
	TopicID          string             `json:"topicID,omitempty"     bson:"topic_id"`
	Style            TemplateStyle      `json:"style,omitempty"       bson:"style"`
	Categories       []TemplateCategory `json:"categories,omitempty"  bson:"categories"`
}

type TemplateCategory struct {
	Name        string           `json:"name,omitempty"        bson:"name"`
	Description string           `json:"description,omitempty" bson:"description"`
	Style       CategoryStyle    `json:"style,omitempty"       bson:"style"`
	Metrics     []TemplateMetric `json:"metrics,omitempty"     bson:"metrics"`
}

type TemplateMetric struct {
	Name  string  `json:"name,omitempty"  bson:"name"`
	Value float64 `json:"value,omitempty" bson:"value"`
	Unit  string  `json:"unit,omitempty"  bson:"unit"`
}

type TemplateStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty"  bson:"icon"`
}

type TemplateTopic struct {
	*base.BaseEntity `                             bson:",inline"`
	Name             string `json:"name,omitempty"        bson:"name"`
	Description      string `json:"description,omitempty" bson:"description"`
}
