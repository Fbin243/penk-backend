package entity

import "tenkhours/pkg/db/base"

type Character struct {
	*base.BaseEntity `                            bson:",inline"`
	ProfileID        string     `json:"profileID,omitempty"  bson:"profile_id"`
	Name             string     `json:"name,omitempty"       bson:"name"`
	Gender           bool       `json:"gender,omitempty"     bson:"gender"`
	Tags             []string   `json:"tags,omitempty"       bson:"tags"`
	Categories       []Category `json:"categories,omitempty" bson:"categories"`
	Metrics          []Metric   `json:"metrics,omitempty"    bson:"metrics"`
	Vision           Vision     `json:"vision,omitempty"     bson:"vision"`
}

type Vision struct {
	Name        string `json:"name,omitempty"        bson:"name"`
	Description string `json:"description,omitempty" bson:"description"`
}

type Category struct {
	ID          string        `json:"id,omitempty"          bson:"_id"`
	Name        string        `json:"name,omitempty"        bson:"name"`
	Description string        `json:"description,omitempty" bson:"description"`
	Style       CategoryStyle `json:"style,omitempty"       bson:"style"`
}

type CategoryStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty"  bson:"icon"`
}

type Metric struct {
	ID         string  `json:"id,omitempty"         bson:"_id"`
	CategoryID *string `json:"categoryID,omitempty" bson:"category_id"`
	Name       string  `json:"name,omitempty"       bson:"name"`
	Value      float64 `json:"value,omitempty"      bson:"value"`
	Unit       string  `json:"unit,omitempty"       bson:"unit"`
}
