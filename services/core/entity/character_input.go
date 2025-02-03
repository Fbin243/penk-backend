package entity

type CharacterInput struct {
	ID            *string             `json:"id,omitempty"`
	Name          string              `json:"name"                    validate:"min=1,max=50"`
	Gender        bool                `json:"gender"`
	Tags          []string            `json:"tags,omitempty"          validate:"tags_valid,dive"`
	CustomMetrics []CustomMetricInput `json:"customMetrics,omitempty" validate:"dive"`
	Vision        *VisionInput        `json:"vision,omitempty"        validate:"omitempty"`
}

type VisionInput struct {
	Name        string  `json:"name"                  validate:"min=1,max=50"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}

type CustomMetricInput struct {
	ID          *string               `json:"id,omitempty"`
	Name        string                `json:"name"                  validate:"min=1,max=50"`
	Description *string               `json:"description,omitempty" validate:"omitempty,max=255"`
	Style       *MetricStyleInput     `json:"style"`
	Properties  []MetricPropertyInput `json:"properties,omitempty"  validate:"dive"`
}

type MetricPropertyInput struct {
	ID    *string            `json:"id,omitempty"`
	Name  string             `json:"name"         validate:"min=1,max=50"`
	Type  MetricPropertyType `json:"type"`
	Value string             `json:"value"        validate:"omitempty,max=50"`
	Unit  string             `json:"unit"         validate:"omitempty,max=50"`
}

type MetricStyleInput struct {
	Color string `json:"color" validate:"hexcolor"`
	Icon  string `json:"icon"`
}
