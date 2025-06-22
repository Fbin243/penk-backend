package rpc

import (
	"context"

	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
)

func (hdl *CoreHandler) UpsertCharacter(ctx context.Context, req *core.CharacterInput) (*core.Character, error) {
	// Map RPC input to entity input
	characterInput, err := Map[core.CharacterInput, entity.CharacterInput](req, append(UnixTimeConverter, MetricConditionConverter...))
	if err != nil {
		return nil, err
	}

	// Call business logic to upsert character
	character, err := hdl.characterBiz.UpsertCharacter(ctx, *characterInput)
	if err != nil {
		return nil, err
	}

	// Map entity to RPC response
	return Map[entity.Character, core.Character](character, append(UnixTimeConverter, MetricConditionConverter...))
}
