package entity

type CharacterInput struct {
	ID   *string `json:"id,omitempty"`
	Name string  `json:"name"         validate:"min=1,max=50"`
}
