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
		users.GetUserStage{
			Metadata: common.Metadata{
				Name: "Create a new user",
			},
		},
		common.SaveToContextStage{
			Key:      common.User,
			JsonPath: "$.data.user",
		},

		users.UpdateUserStage{
			Metadata: common.Metadata{
				Name: "Update user info",
			},
		},

		characters.CreateCharacterStage{
			Metadata: common.Metadata{
				Name: "Create the first character",
			},
		},
		common.SaveToContextStage{
			Key:      common.Character1,
			JsonPath: "$.data.createCharacter",
		},

		characters.CreateCharacterStage{
			Metadata: common.Metadata{
				Name: "Create the second character",
			},
		},
		common.SaveToContextStage{
			Key:      common.Character2,
			JsonPath: "$.data.createCharacter",
		},

		characters.CreateCharacterStage{
			Metadata: common.Metadata{
				Name:        "Create the third character",
				ExpectError: true,
			},
		},

		characters.UpdateCharacterStage{
			Metadata: common.Metadata{
				Name: "Update info of the first character",
			},
			CharacterKey: common.Character1,
		},

		metrics.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Name: "Create the first metric",
			},
			CharacterKey: common.Character1,
		},

		metrics.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Name: "Create the second metric",
			},

			CharacterKey: common.Character1,
		},

		metrics.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Name:        "Create the third metric",
				ExpectError: false,
			},

			CharacterKey: common.Character1,
		},

		common.SwitchUrlStage{
			NewUrl: common.AnalyticsUrl,
		},
	)

	err := p(&ctx)
	if err != nil {
		common.LogResponse()
	}

	assert.Empty(t, err)
}
