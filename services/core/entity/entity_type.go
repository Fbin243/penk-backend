package entity

type EntityType string

const (
	EntityTypeCategory EntityType = "Category"
	EntityTypeMetric   EntityType = "Metric"
	EntityTypeGoal     EntityType = "Goal"
	EntityTypeHabit    EntityType = "Habit"
	EntityTypeTask     EntityType = "Task"
)
