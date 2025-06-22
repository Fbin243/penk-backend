package errors

import (
	"strings"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	ErrUnauthorized     = NewGQLError(ErrCodeUnauthorized, "unauthorized")
	ErrPermissionDenied = NewGQLError(ErrCodePermissionDenied, "permission denied")
	ErrBadRequest       = NewGQLError(ErrCodeBadRequest, "bad request")
	ErrNotFound         = NewGQLError(ErrCodeNotFound, "not found")
	ErrMongoNotFound    = NewGQLError(ErrCodeMongoNotFound, "mongo not found")
	ErrRedisNotFound    = NewGQLError(ErrCodeRedisNotFound, "redis not found")
	ErrLimitMetric      = NewGQLError(ErrCodeLimitMetric, "over limit metric")
	ErrLimitCheckbox    = NewGQLError(ErrCodeLimitCheckbox, "over limit checkbox")
	ErrLimitCharacter   = NewGQLError(ErrCodeLimitCharacter, "over limit character")
	ErrLimitCategory    = NewGQLError(ErrCodeLimitCategory, "over limit category")
	ErrLimitGoal        = NewGQLError(ErrCodeLimitGoal, "over limit goal")
	ErrLimitHabit       = NewGQLError(ErrCodeLimitHabit, "over limit habit")
	ErrLimitTask        = NewGQLError(ErrCodeLimitTask, "over limit task")
)

func NewGQLError(code ErrorCode, msg any) *gqlerror.Error {
	var message string

	switch v := msg.(type) {
	case error:
		message = v.Error()
	case string:
		message = v
	default:
		message = CodeToMessage(code)
	}

	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			"code": code,
		},
	}
}

func CodeToMessage(code ErrorCode) string {
	return strings.ToLower(strings.Join(strings.Split(string(code), "_"), " "))
}

func HasCode(err error, code ErrorCode) bool {
	if gqlErr, ok := err.(*gqlerror.Error); ok {
		if gqlErr != nil {
			return gqlErr.Extensions["code"] == code
		}
	}

	return false
}
