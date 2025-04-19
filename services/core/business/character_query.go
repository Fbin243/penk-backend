package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/services/core/entity"
)

func (biz *CharacterBusiness) GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	characters, err := biz.CharacterRepo.GetCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	return characters, nil
}
