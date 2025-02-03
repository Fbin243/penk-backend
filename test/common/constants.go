package common

type ContextKey int

const (
	TestingT ContextKey = iota
	Profile
	Character
	Snapshot
	TimeTracking
)

type Case int

const (
	GetProfile Case = iota + 1
	UpdateProfile
	CreateCharacter
	UpdateCharacter
	CreateMetrics
	UpdateMetrics
	DeleteMetrics
	CreateProperties
	UpdateProperties
	DeleteProperties
	TimeTrackingWithoutMetric
	TimeTrackingWithMetric
)

type FlowKey string

const (
	ProfileFlowKey      FlowKey = "profile"
	CharacterFlowKey    FlowKey = "character"
	TimeTrackingFlowKey FlowKey = "timetracking"
	SnapshotFlowKey     FlowKey = "snapshot"
)
