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

type UpdateCustomMetricCase int

const (
	CreateProperties UpdateCustomMetricCase = iota
	UpdateProperties
	DeleteProperties
)
