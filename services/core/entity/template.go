package entity

import "tenkhours/pkg/db/base"

type Template struct {
	*base.BaseEntity `                             bson:",inline"`
	Name             string           `json:"name,omitempty"        bson:"name"`
	Description      string           `json:"description,omitempty" bson:"description"`
	CategoryID       string           `json:"categoryID,omitempty"  bson:"category_id"`
	Style            TemplateStyle    `json:"style,omitempty"       bson:"style"`
	Metrics          []TemplateMetric `json:"metrics,omitempty"     bson:"metrics"`
}

type TemplateMetric struct {
	Name        string             `json:"name,omitempty"        bson:"name"`
	Description string             `json:"description,omitempty" bson:"description"`
	Style       MetricStyle        `json:"style,omitempty"       bson:"style"`
	Properties  []TemplateProperty `json:"properties,omitempty"  bson:"properties"`
}

type TemplateProperty struct {
	Name  string             `json:"name,omitempty"  bson:"name"`
	Type  MetricPropertyType `json:"type,omitempty"  bson:"type"`
	Value string             `json:"value,omitempty" bson:"value"`
	Unit  string             `json:"unit,omitempty"  bson:"unit"`
}

type TemplateStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty"  bson:"icon"`
}

type TemplateCategory struct {
	*base.BaseEntity `                             bson:",inline"`
	Name             string `json:"name,omitempty"        bson:"name"`
	Description      string `json:"description,omitempty" bson:"description"`
}
