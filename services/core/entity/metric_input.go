package entity

type MetricInput struct {
	ID          *string `json:"id,omitempty"`
	CharacterID string  `json:"characterID"`
	CategoryID  *string `json:"categoryID,omitempty"`
	Name        string  `json:"name"                 validate:"min=1,max=50"`
	Value       float64 `json:"value"`
	Unit        string  `json:"unit"                 validate:"omitempty,max=50"`
}
