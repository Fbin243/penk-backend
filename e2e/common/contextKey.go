package common

type ContextKey string

const (
	TestingT           ContextKey = "TestingT"
	Profile            ContextKey = "Profile"
	AnotherCharacter   ContextKey = "AnotherCharacter"
	CurrentCharacter   ContextKey = "CurrentCharacter"
	Snapshot           ContextKey = "Snapshot"
	TimeTracking       ContextKey = "TimeTracking"
	FirstCustomMetric  ContextKey = "FirstCustomMetric"
	SecondCustomMetric ContextKey = "SecondCustomMetric"
	MetricProperty     ContextKey = "MetricProperty"
)
