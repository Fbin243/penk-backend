package errors

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func DefaultPresenter(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)

	if err.Message == "" {
		err.Message = "internal server error"
	}

	if err.Extensions["code"] == nil {
		err.Extensions = map[string]interface{}{
			"code": ErrCodeInternalServer,
		}
	}

	return err
}
