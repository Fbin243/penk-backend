package entity

type ProfileInput struct {
	Name               string `json:"name"                         validate:"min=1,max=50"`
	ImageURL           string `json:"imageURL"`
	CurrentCharacterID string `json:"currentCharacterID,omitempty"`
}
