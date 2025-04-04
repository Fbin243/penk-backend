package entity

type HabitInput struct {
	ID             *string        `json:"id,omitempty"`
	CategoryID     *string        `json:"categoryID"`
	CompletionType CompletionType `json:"completionType"`
	Name           string         `json:"name"`
	Value          float64        `json:"value"`
	Unit           *string        `json:"unit"`
	RRule          string         `json:"rrule"`
}
