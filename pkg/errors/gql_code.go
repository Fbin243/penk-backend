package errors

type ErrorCode string

const (
	// General errors
	ErrCodeInternalServer   ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest       ErrorCode = "BAD_REQUEST"
	ErrCodePermissionDenied ErrorCode = "PERMISSION_DENIED"
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"

	// Core errors
	ErrCodeLimitCharacter      ErrorCode = "LIMIT_CHARACTER"
	ErrCodeLimitMetric         ErrorCode = "LIMIT_METRIC"
	ErrCodeLimitProperty       ErrorCode = "LIMIT_PROPERTY"
	ErrCodeGoalAlreadyFinished ErrorCode = "GOAL_ALREADY_FINISHED"
	ErrCodeGoalAlreadyExpired  ErrorCode = "GOAL_ALREADY_EXPIRED"

	// Time tracking errors
	ErrCodeUnderMinDuration          ErrorCode = "UNDER_MIN_DURATION"
	ErrCodeOverMaxDifferenceDuration ErrorCode = "OVER_MAX_DIFFERENCE_DURATION"
	ErrCodeTimeTrackingAlreadyExists ErrorCode = "TIME_TRACKING_ALREADY_EXISTS"

	// Analytics errors
	ErrCodeLimitSnapshot     ErrorCode = "LIMIT_SNAPSHOT"
	ErrCodeDuplicateSnapshot ErrorCode = "DUPLICATE_SNAPSHOT"
)
