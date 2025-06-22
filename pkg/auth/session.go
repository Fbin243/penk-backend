package auth

import (
	"context"

	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
)

func GetAuthSession(ctx context.Context) (rdb.AuthSession, error) {
	authSession, ok := ctx.Value(AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return rdb.AuthSession{}, errors.ErrUnauthorized
	}

	return authSession, nil
}
