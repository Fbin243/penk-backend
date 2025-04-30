package common

type ContextKey int

const (
	TestingT ContextKey = iota
	Profile
	Character
	Snapshot
	TimeTracking
	Goal
)

type Case int

const (
	GetProfile Case = iota + 1
	UpdateProfile
	CreateCharacter
	UpdateCharacter
	CreateCategories
	UpdateCategories
	DeleteCategories
	CreateMetrics
	UpdateMetrics
	DeleteMetrics
	TimeTrackingWithoutCategory
	TimeTrackingWithCategory
	CreateGoal
	UpdateGoal
	DeleteGoal
	CurrentTimeTrackingExist
)

type FlowKey string

const (
	ProfileFlowKey      FlowKey = "profile"
	CharacterFlowKey    FlowKey = "character"
	TimeTrackingFlowKey FlowKey = "timetracking"
	SnapshotFlowKey     FlowKey = "snapshot"
	GoalFlowKey         FlowKey = "goal"
)
