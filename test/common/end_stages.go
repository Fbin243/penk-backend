package common

type EndStage string

const (
	// Key for context
	EndStageKey EndStage = "END_STAGE"
	// End stages
	ProfileStage      EndStage = "profile"
	CharacterStage    EndStage = "character"
	CustomMetricStage EndStage = "custom_metric"
)
