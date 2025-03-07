package rpc

import (
	"context"

	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

func (hdl *CoreHandler) UpsertCharacter(ctx context.Context, req *core.CharacterInput) (*core.Character, error) {
	resp := &core.Character{}

	characterInput := entity.CharacterInput{}
	copier.Copy(&characterInput, req)
	character, err := hdl.characterBiz.UpsertCharacter(ctx, characterInput)
	if err != nil {
		return resp, err
	}

	copier.Copy(resp, character)
	resp.CreatedAt = character.CreatedAt.Unix()
	resp.UpdatedAt = character.UpdatedAt.Unix()

	return resp, nil
}
