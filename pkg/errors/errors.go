package errors

import "github.com/vektah/gqlparser/v2/gqlerror"

func NewError(code ErrorCode, msg any) *gqlerror.Error {
	var message string
	switch v := msg.(type) {
	case error:
		message = v.Error()
	}

	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			"code": code,
		},
	}
}

func Unauthorized() *gqlerror.Error {
	return NewError(ErrCodeUnauthorized, "unauthorized")
}

func PermissionDenied() *gqlerror.Error {
	return NewError(ErrCodePermissionDenied, "permission denied")
}
