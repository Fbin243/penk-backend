package test

import (
	"context"
	"testing"

	"tenkhours/pineline"

	"github.com/stretchr/testify/assert"
)

func TestUserFlow(t *testing.T) {
	ctx := context.WithValue(context.Background(), TestingT, t)

	p := pineline.Pineline(
		getUser(false),
		saveToContext(User, "$.data.user"),
		logResponse,

		updateUser(false),
		logResponse,

		createCharacter(false),
		saveToContext(Character1, "$.data.createCharacter"),
		logResponse,

		createCharacter(false),
		saveToContext(Character2, "$.data.createCharacter"),
		logResponse,

		createCharacter(true),
		logResponse,

		updateCharacter(Character1, false),
		logResponse,

		createCustomMetric(Character1, false),
		logResponse,

		createCustomMetric(Character1, false),
		logResponse,

		createCustomMetric(Character1, true),
		logResponse,

		deleteCharacter(Character2, false),
		logResponse,
	)

	err := p(&ctx)
	if err != nil {
		logResponse(&ctx)
	}

	assert.Empty(t, err)
}
