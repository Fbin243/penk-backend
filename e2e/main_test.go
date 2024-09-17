package test

import (
	"context"
	"testing"

	"tenkhours/pineline"
	analytics_characters "tenkhours/test/analytics/characters"
	"tenkhours/test/common"
	"tenkhours/test/core/characters"
	"tenkhours/test/core/characters/metrics"
	"tenkhours/test/core/characters/metrics/properties"
	"tenkhours/test/core/timetrackings"
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

		// --------- CHARACTER -----------

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

		// --------- TIME TRACKING WITHOUT METRIC -----------

		timetrackings.CreateTimeTrackingStage{
			Metadata: common.Metadata{
				Describe: "Start a new tracking session of the current character",
			},
			CharacterKey: common.CurrentCharacter,
		},

		common.SaveToContextStage{
			Key:      common.TimeTracking,
			JsonPath: "data.createTimeTracking",
		},

		timetrackings.CreateTimeTrackingStage{
			Metadata: common.Metadata{
				Describe:    "Start an existing session again",
				ExpectError: true,
			},
			CharacterKey: common.CurrentCharacter,
		},

		timetrackings.UpdateTimeTracking{
			Metadata: common.Metadata{
				Describe: "End an existing tracking session of the current character",
			},
			TimeTrackingKey: common.TimeTracking,
		},

		timetrackings.UpdateTimeTracking{
			Metadata: common.Metadata{
				Describe:    "End an finished session again",
				ExpectError: true,
			},
			TimeTrackingKey: common.TimeTracking,
		},

		characters.UpdateCharacterStage{
			Metadata: common.Metadata{
				Describe: "Update info of the current character",
			},
			CharacterKey: common.CurrentCharacter,
		},

		// --------- CUSTOM METRIC -----------

		metrics.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Create the first metric",
			},
			CharacterKey: common.CurrentCharacter,
		},

		common.SaveToContextStage{
			Key:      common.FirstCustomMetric,
			JsonPath: "data.createCustomMetric",
		},

		metrics.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Create the second metric",
			},
			CharacterKey: common.CurrentCharacter,
		},

		common.SaveToContextStage{
			Key:      common.SecondCustomMetric,
			JsonPath: "data.createCustomMetric",
		},

		metrics.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe:    "Create the third metric",
				ExpectError: true,
			},
			CharacterKey: common.CurrentCharacter,
		},

		metrics.UpdateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Update info for the first custom metric",
			},
			CustomMetricKey: common.FirstCustomMetric,
			CharacterKey:    common.CurrentCharacter,
		},

		metrics.ResetCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Reset the second custom metric",
			},
			CustomMetricKey: common.SecondCustomMetric,
			CharacterKey:    common.CurrentCharacter,
		},

		// --------- TIME TRACKING WITH METRIC -----------

		timetrackings.CreateTimeTrackingStage{
			Metadata: common.Metadata{
				Describe: "Start a new tracking session of the current character with the first custom metric",
			},
			CharacterKey:    common.CurrentCharacter,
			CustomMetricKey: common.FirstCustomMetric,
			TrackWithMetric: true,
		},

		common.SaveToContextStage{
			Key:      common.TimeTracking,
			JsonPath: "data.createTimeTracking",
		},

		timetrackings.UpdateTimeTracking{
			Metadata: common.Metadata{
				Describe: "End an existing tracking session of the current character with the first custom metric",
			},
			CustomMetricKey: common.FirstCustomMetric,
			TimeTrackingKey: common.TimeTracking,
			TrackWithMetric: true,
		},

		// --------- CUSTOM METRIC PROPERTY -----------

		properties.CreateMetricPropertyStage{
			Metadata: common.Metadata{
				Describe: "Create the first property for the second metric",
			},
			CharacterKey:    common.CurrentCharacter,
			CustomMetricKey: common.SecondCustomMetric,
		},

		properties.CreateMetricPropertyStage{
			Metadata: common.Metadata{
				Describe: "Create the second property for the second metric",
			},
			CharacterKey:    common.CurrentCharacter,
			CustomMetricKey: common.SecondCustomMetric,
		},

		common.SaveToContextStage{
			Key:      common.MetricProperty,
			JsonPath: "data.createMetricProperty",
		},

		properties.CreateMetricPropertyStage{
			Metadata: common.Metadata{
				Describe:    "Create the third property for the second metric",
				ExpectError: true,
			},
			CharacterKey:    common.CurrentCharacter,
			CustomMetricKey: common.SecondCustomMetric,
		},

		properties.UpdateMetricPropertyStage{
			Metadata: common.Metadata{
				Describe: "Update the second property of the second metric",
			},
			CharacterKey:      common.CurrentCharacter,
			CustomMetricKey:   common.SecondCustomMetric,
			MetricPropertyKey: common.MetricProperty,
		},

		properties.DeleteMetricPropertyStage{
			Metadata: common.Metadata{
				Describe: "Delete the second property of the second metric",
			},
			CharacterKey:      common.CurrentCharacter,
			CustomMetricKey:   common.SecondCustomMetric,
			MetricPropertyKey: common.MetricProperty,
		},

		metrics.DeleteCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Delete the second custom metric",
			},
			CustomMetricKey: common.SecondCustomMetric,
			CharacterKey:    common.CurrentCharacter,
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

		// --------- SNAPSHOT -----------

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

		characters.DeleteCharacterStage{
			Metadata: common.Metadata{
				Describe: "Delete another character",
			},
			CharacterKey: common.AnotherCharacter,
		},

		users.GetUserStage{
			Metadata: common.Metadata{
				Describe: "Get updated user info for reviewing after performing flow",
			},
		},
	)

	err := p(&ctx)
	if err != nil {
		common.LogResponse()
	}

	assert.Empty(t, err)
}
