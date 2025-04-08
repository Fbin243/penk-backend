package errors

type ErrorCode string

const (
	// General errors
	ErrCodeInternalServer   ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest       ErrorCode = "BAD_REQUEST"
	ErrCodePermissionDenied ErrorCode = "PERMISSION_DENIED"
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeMongoNotFound    ErrorCode = "MONGO_NOT_FOUND"
	ErrCodeRedisNotFound    ErrorCode = "REDIS_NOT_FOUND"

	// Core errors
	ErrCodeLimitCharacter            ErrorCode = "LIMIT_CHARACTER"
	ErrCodeLimitCategory             ErrorCode = "LIMIT_CATEGORY"
	ErrCodeLimitGoal                 ErrorCode = "LIMIT_GOAL"
	ErrCodeLimitMetric               ErrorCode = "LIMIT_METRIC"
	ErrCodeLimitCheckbox             ErrorCode = "LIMIT_CHECKBOX"
	ErrCodeLimitHabit                ErrorCode = "LIMIT_HABIT"
	ErrCodeLimitTask                 ErrorCode = "LIMIT_TASK"
	ErrCodeUnderMinDuration          ErrorCode = "UNDER_MIN_DURATION"
	ErrCodeOverMaxDifferenceDuration ErrorCode = "OVER_MAX_DIFFERENCE_DURATION"
	ErrCodeTimeTrackingAlreadyExists ErrorCode = "TIME_TRACKING_ALREADY_EXISTS"
)
