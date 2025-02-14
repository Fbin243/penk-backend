package entity

type CharacterInput struct {
	ID         *string         `json:"id,omitempty"`
	Name       string          `json:"name"                 validate:"min=1,max=50"`
	Gender     bool            `json:"gender"`
	Tags       []string        `json:"tags,omitempty"       validate:"tags_valid,dive"`
	Categories []CategoryInput `json:"categories,omitempty" validate:"dive"`
	Metrics    []MetricInput   `json:"metrics,omitempty"    validate:"dive"`
	Vision     *VisionInput    `json:"vision,omitempty"     validate:"omitempty"`
}

type VisionInput struct {
	Name        string  `json:"name"                  validate:"min=1,max=50"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}

type CategoryInput struct {
	ID          *string             `json:"id,omitempty"`
	Name        string              `json:"name"                  validate:"min=1,max=50"`
	Description *string             `json:"description,omitempty" validate:"omitempty,max=255"`
	Style       *CategoryStyleInput `json:"style"`
}

type MetricInput struct {
	ID         *string `json:"id,omitempty"`
	CategoryID *string `json:"category_id,omitempty"`
	Name       string  `json:"name"                  validate:"min=1,max=50"`
	Value      float64 `json:"value"`
	Unit       string  `json:"unit"                  validate:"omitempty,max=50"`
}

type CategoryStyleInput struct {
	Color string `json:"color" validate:"hexcolor"`
	Icon  string `json:"icon"`
}
