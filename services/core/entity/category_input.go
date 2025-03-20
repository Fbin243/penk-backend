package entity

type CategoryInput struct {
	ID          *string             `json:"id,omitempty"`
	CharacterID string              `json:"characterID,omitempty"`
	Name        string              `json:"name"                  validate:"min=1,max=50"`
	Description *string             `json:"description,omitempty" validate:"omitempty,max=255"`
	Style       *CategoryStyleInput `json:"style"`
}

type CategoryStyleInput struct {
	Color string `json:"color" validate:"hexcolor"`
	Icon  string `json:"icon"`
}
