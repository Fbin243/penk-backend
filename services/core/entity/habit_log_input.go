package entity

type HabitLogInput struct {
	Timestamp string  `json:"timestamp"`
	HabitID   string  `json:"habitID"`
	Value     float64 `json:"value"`
}
