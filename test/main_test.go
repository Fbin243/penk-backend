package test

import (
	"context"
	"testing"

	"tenkhours/pineline"
	"tenkhours/test/common"
	"tenkhours/test/core/characters"
	"tenkhours/test/core/characters/metrics"
	"tenkhours/test/core/users"

	"github.com/stretchr/testify/assert"
)

func TestUserFlow(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.TestingT, t)

	p := pineline.Pineline(
		users.GetUser(false),
		common.SaveToContext(common.User, "$.data.user"),
		common.LogResponse,

		users.UpdateUser(false),
		common.LogResponse,

		characters.CreateCharacter(false),
		common.SaveToContext(common.Character1, "$.data.createCharacter"),
		common.LogResponse,

		characters.CreateCharacter(false),
		common.SaveToContext(common.Character2, "$.data.createCharacter"),
		common.LogResponse,

		characters.CreateCharacter(true),
		common.LogResponse,

		characters.UpdateCharacter(common.Character1, false),
		common.LogResponse,

		metrics.CreateCustomMetric(common.Character1, false),
		common.LogResponse,

		metrics.CreateCustomMetric(common.Character1, false),
		common.LogResponse,

		metrics.CreateCustomMetric(common.Character1, true),
		common.LogResponse,

		characters.DeleteCharacter(common.Character2, false),
		common.LogResponse,
	)

	err := p(&ctx)
	if err != nil {
		common.LogResponse(&ctx)
	}

	assert.Empty(t, err)
}
