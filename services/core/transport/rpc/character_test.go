package rpc_test

import (
	"testing"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var character = &entity.Character{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	},
	ProfileID: mongodb.GenObjectID(),
	Name:      "Character name",
}

var category = entity.Category{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	},
	CharacterID: character.ID,
	Name:        "Category name",
	Description: "Category desc",
	Style: entity.CategoryStyle{
		Color: "#000000",
		Icon:  "icon.png",
	},
}

var metric = entity.Metric{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	},
	CharacterID: character.ID,
	CategoryID:  lo.ToPtr(category.ID),
	Name:        "Metric name",
	Value:       1.0,
	Unit:        "unit",
}

func TestMapCharacter(t *testing.T) {
	rpcCharacter := &core.Character{}
	copier.Copy(rpcCharacter, character)
	rpcCharacter.CreatedAt = character.CreatedAt.Unix()
	rpcCharacter.UpdatedAt = character.UpdatedAt.Unix()

	assert.Equal(t, character.ID, rpcCharacter.Id)
	assert.Equal(t, character.CreatedAt.Unix(), rpcCharacter.CreatedAt)
	assert.Equal(t, character.UpdatedAt.Unix(), rpcCharacter.UpdatedAt)
	assert.Equal(t, character.ProfileID, rpcCharacter.ProfileId)
	assert.Equal(t, character.Name, rpcCharacter.Name)
}

var rpcCharacterInput = &core.CharacterInput{
	Id:     lo.ToPtr(mongodb.GenObjectID()),
	Name:   "Example",
	Gender: false,
	Tags:   []string{"#tag1", "#tag2"},
}

func TestMapCharacterInput(t *testing.T) {
	entityCharacterInput := &entity.CharacterInput{}
	copier.Copy(entityCharacterInput, rpcCharacterInput)

	assert.Equal(t, entityCharacterInput.ID, rpcCharacterInput.Id)
	assert.Equal(t, entityCharacterInput.Name, rpcCharacterInput.Name)
}
