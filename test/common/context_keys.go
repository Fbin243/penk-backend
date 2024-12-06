package common

type ContextKey int

const (
	TestingT ContextKey = iota
	Profile
	AnotherCharacter
	CurrentCharacter
	Snapshot
	TimeTracking
	FirstCustomMetric
	SecondCustomMetric
	MetricProperty
	MetricProperties
)

type (
	// Character
	CreateCharacterCase int

	// Custom metric
	CreateCustomMetricCase int
	UpdateCustomMetricCase int
)

const (
	// Character cases
	CreateCharacterWithCustomMetrics CreateCharacterCase = iota + 1

	// Custom metric cases
	CreateMetricWithProperties CreateCustomMetricCase = iota + 1
	CreateProperties           UpdateCustomMetricCase = iota + 1
	UpdateProperties
	DeleteProperties
)
