package entity

import "time"

type HabitInput struct {
	ID             *string        `json:"id,omitempty"`
	CharacterID    string         `json:"characterID"`
	CategoryID     *string        `json:"categoryID"`
	CompletionType CompletionType `json:"completionType"`
	Name           string         `json:"name"`
	Value          float64        `json:"value"`
	Unit           *string        `json:"unit"`
	StartTime      time.Time      `json:"startTime"`
	EndTime        *time.Time     `json:"endTime"`
	Frequency      string         `json:"frequency"`
}
