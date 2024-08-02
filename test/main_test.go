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
		getUser(),
		saveToContext(UserID, "$.data.user.id"),
		logResponse,

		updateUser(),
		logResponse,

		createNewCharacter(false),
		saveToContext(CharacterID, "$.data.createCharacter.id"),
		logResponse,

		createNewCharacter(false),
		logResponse,

		createNewCharacter(true),
		logResponse,

		updateCharacter(),
		logResponse,

		createCustomMetric(false),
		logResponse,

		createCustomMetric(false),
		logResponse,

		createCustomMetric(true),
		logResponse,

		deleteCharacter(),
		logResponse,
	)

	err := p(&ctx)
	assert.Empty(t, err)
	logResponse(&ctx)
}
