package common

type ContextKey string

const (
	TestingT           ContextKey = "TestingT"
	User               ContextKey = "User"
	AnotherCharacter   ContextKey = "AnotherCharacter"
	CurrentCharacter   ContextKey = "CurrentCharacter"
	Snapshot           ContextKey = "Snapshot"
	TimeTracking       ContextKey = "TimeTracking"
	FirstCustomMetric  ContextKey = "FirstCustomMetric"
	SecondCustomMetric ContextKey = "SecondCustomMetric"
	MetricProperty     ContextKey = "MetricProperty"
)
