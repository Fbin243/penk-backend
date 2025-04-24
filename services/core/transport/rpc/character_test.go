package rpc_test

import (
	"testing"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
	"tenkhours/services/core/transport/rpc"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// Entity
var character = entity.Character{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	},
	ProfileID: mongodb.GenObjectID(),
	Name:      "Character name",
}

// RPC Input
var rpcCharacterInput = &core.CharacterInput{
	Id:   lo.ToPtr(mongodb.GenObjectID()),
	Name: "Character name",
}

func TestMapCharacter(t *testing.T) {
	rpcCharacter, err := rpc.Map[entity.Character, core.Character](&character, rpc.UnixTimeConverter)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, character.ID, rpcCharacter.Id)
	assert.Equal(t, character.CreatedAt.Unix(), rpcCharacter.CreatedAt)
	assert.Equal(t, character.UpdatedAt.Unix(), rpcCharacter.UpdatedAt)
	assert.Equal(t, character.ProfileID, rpcCharacter.ProfileId)
	assert.Equal(t, character.Name, rpcCharacter.Name)
}

func TestMapCharacterInput(t *testing.T) {
	characterInput, err := rpc.Map[core.CharacterInput, entity.CharacterInput](rpcCharacterInput, rpc.UnixTimeConverter)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, rpcCharacterInput.Id, characterInput.ID)
	assert.Equal(t, rpcCharacterInput.Name, characterInput.Name)
}
