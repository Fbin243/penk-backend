package test

import (
	"context"
	"testing"

	"tenkhours/pineline"
	analytics_characters "tenkhours/test/analytics/characters"
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
				Describe: "Create a new user",
			},
			CreateNewUser: true,
		},

		common.SaveToContextStage{
			Key:      common.User,
			JsonPath: "data.user",
		},

		users.UpdateUserStage{
			Metadata: common.Metadata{
				Describe: "Update user info",
			},
		},

		characters.CreateCharacterStage{
			Metadata: common.Metadata{
				Describe: "Create the first character",
			},
		},

		common.SaveToContextStage{
			Key:      common.AnotherCharacter,
			JsonPath: "data.createCharacter",
		},

		characters.CreateCharacterStage{
			Metadata: common.Metadata{
				Describe: "Create the second character",
			},
		},

		common.SaveToContextStage{
			Key:      common.CurrentCharacter,
			JsonPath: "data.createCharacter",
		},

		characters.CreateCharacterStage{
			Metadata: common.Metadata{
				Describe:    "Create the third character",
				ExpectError: true,
			},
		},

		characters.UpdateCharacterStage{
			Metadata: common.Metadata{
				Describe: "Update info of the current character",
			},
			CharacterKey: common.CurrentCharacter,
		},

		metrics.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Create the first metric",
			},
			CharacterKey: common.CurrentCharacter,
		},

		metrics.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Create the second metric",
			},
			CharacterKey: common.CurrentCharacter,
		},

		metrics.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe:    "Create the third metric",
				ExpectError: true,
			},
			CharacterKey: common.CurrentCharacter,
		},

		users.GetUserStage{
			Metadata: common.Metadata{
				Describe: "Get updated user info after add characters, metrics, etc..",
			},
		},

		common.SaveToContextStage{
			Key:      common.User,
			JsonPath: "data.user",
		},

		analytics_characters.CreateSnapshotStage{
			Metadata: common.Metadata{
				Describe: "Create a snapshot for the current character",
			},
			CharacterKey: common.CurrentCharacter,
		},

		common.SaveToContextStage{
			Key:      common.Snapshot,
			JsonPath: "data.createSnapshot",
		},

		analytics_characters.GetCharacterSnapshotsStage{
			Metadata: common.Metadata{
				Describe: "Get all snapshots of the current character after add a new snapshot for it",
			},
			SnapshotKey:    common.Snapshot,
			HasOneSnapshot: true,
		},

		analytics_characters.CreateSnapshotStage{
			Metadata: common.Metadata{
				Describe: "Create a snapshot for another character",
			},
			CharacterKey: common.AnotherCharacter,
		},

		analytics_characters.GetUserSnapshotsStage{
			Metadata: common.Metadata{
				Describe: "Get all snapshots of the user after add a new snapshot for each of characters, so we have 2 snapshots",
			},
			HasTwoSnapshots: true,
		},
	)

	err := p(&ctx)
	if err != nil {
		common.LogResponse()
	}

	assert.Empty(t, err)
}
