package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

func (biz *CharacterBusiness) GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	characters, err := biz.CharacterRepo.GetCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	return characters, nil
}
