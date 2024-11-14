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
	UpdateCustomMetricCase int
	CreateCustomMetricCase int
)

const (
	CreateProperties UpdateCustomMetricCase = iota + 1
	UpdateProperties
	DeleteProperties
	CreateMetricWithProperties CreateCustomMetricCase = iota + 1
)
