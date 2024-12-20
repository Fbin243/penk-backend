package errors

import (
	"strings"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

func NewError(code ErrorCode, msg any) *gqlerror.Error {
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

func Unauthorized() *gqlerror.Error {
	return NewError(ErrCodeUnauthorized, "unauthorized")
}

func PermissionDenied() *gqlerror.Error {
	return NewError(ErrCodePermissionDenied, "permission denied")
}
