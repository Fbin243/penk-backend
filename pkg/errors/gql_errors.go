package errors

import (
	"strings"

	"github.com/vektah/gqlparser/v2/gqlerror"
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

func Unauthorized() *gqlerror.Error {
	return NewGQLError(ErrCodeUnauthorized, "unauthorized")
}

func PermissionDenied() *gqlerror.Error {
	return NewGQLError(ErrCodePermissionDenied, "permission denied")
}
