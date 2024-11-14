package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"tenkhours/pkg/pineline"
	"tenkhours/test/analytics"
	"tenkhours/test/common"
	"tenkhours/test/core"

	"tenkhours/test/timetrackings"

	"github.com/stretchr/testify/assert"
)

func TestUserFlow(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.TestingT, t)

	token := os.Getenv("TOKEN")
	if token == "" {
		t.Fatalf("Token not set in environment variables")
	}

	fmt.Printf("Token: %v\n", token)
	common.IdToken = token

	// Get the end stage of the test flow and put it into the context
	endStage := os.Getenv("END_STAGE")
	fmt.Printf("End stage: %v\n", endStage)
	ctx = context.WithValue(ctx, common.EndStageKey, endStage)

	p := pineline.Pineline(

		// --------- START >> PROFILE -----------

		core.GetProfileStage{
			Metadata: common.Metadata{
				Describe: "Create a new profile",
			},
			CreateNewProfile: true,
		},

		common.SaveToContextStage{
			Key:      common.Profile,
			JsonPath: "data.profile",
		},

		core.UpdateProfileStage{
			Metadata: common.Metadata{
				Describe: "Update profile info",
			},
		},

		common.SaveToContextStage{
			Key:      common.Profile,
			JsonPath: "data.updateProfile",
		},

		common.CheckEndStage{
			CurrentStage: common.ProfileStage,
		},

		// --------- END << PROFILE -----------

		// --------- START >> CHARACTER -----------

		core.CreateCharacterStage{
			Metadata: common.Metadata{
				Describe: "Create the first character",
			},
			CreateCharacterCase: common.CreateCharacterWithCustomMetrics,
		},

		common.SaveToContextStage{
			Key:      common.AnotherCharacter,
			JsonPath: "data.createCharacter",
		},

		core.CreateCharacterStage{
			Metadata: common.Metadata{
				Describe: "Create the second character",
			},
		},

		common.SaveToContextStage{
			Key:      common.CurrentCharacter,
			JsonPath: "data.createCharacter",
		},

		core.CreateCharacterStage{
			Metadata: common.Metadata{
				Describe:    "Create the third character",
				ExpectError: true,
			},
		},

		core.UpdateCharacterStage{
			Metadata: common.Metadata{
				Describe: "Update info of the current character",
			},
			CharacterKey: common.CurrentCharacter,
		},

		common.CheckEndStage{
			CurrentStage: common.CharacterStage,
		},

		// --------- END << CHARACTER -----------

		// --------- START >> CUSTOM METRIC -----------
		core.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Create the first metric",
			},
			CharacterKey: common.CurrentCharacter,
		},

		common.SaveToContextStage{
			Key:      common.FirstCustomMetric,
			JsonPath: "data.createCustomMetric",
		},

		core.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Create the second metric with properties",
			},
			CharacterKey: common.CurrentCharacter,
			Case:         common.CreateMetricWithProperties,
		},

		common.SaveToContextStage{
			Key:      common.SecondCustomMetric,
			JsonPath: "data.createCustomMetric",
		},

		core.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Create the third metric",
			},
			CharacterKey: common.CurrentCharacter,
		},

		core.CreateCustomMetricStage{
			Metadata: common.Metadata{
				Describe:    "Create the fourth metric",
				ExpectError: true,
			},
			CharacterKey: common.CurrentCharacter,
		},

		common.CheckEndStage{
			CurrentStage: common.CustomMetricStage,
		},

		// --------- END << CUSTOM METRIC -----------

		// --------- TIME TRACKING WITHOUT METRIC -----------

		// timetrackings.CreateTimeTrackingStage{
		// 	Metadata: common.Metadata{
		// 		Describe: "Start a new tracking session of the current character",
		// 	},
		// 	CharacterKey: common.CurrentCharacter,
		// },

		// common.SaveToContextStage{
		// 	Key:      common.TimeTracking,
		// 	JsonPath: "data.createTimeTracking",
		// },

		// timetrackings.CreateTimeTrackingStage{
		// 	Metadata: common.Metadata{
		// 		Describe:    "Start an existing session again",
		// 		ExpectError: true,
		// 	},
		// 	CharacterKey: common.CurrentCharacter,
		// },

		// timetrackings.UpdateTimeTracking{
		// 	Metadata: common.Metadata{
		// 		Describe: "End an existing tracking session of the current character",
		// 	},
		// 	TimeTrackingKey: common.TimeTracking,
		// },

		// timetrackings.UpdateTimeTracking{
		// 	Metadata: common.Metadata{
		// 		Describe:    "End an finished session again",
		// 		ExpectError: true,
		// 	},
		// 	TimeTrackingKey: common.TimeTracking,
		// },

		// --------- CUSTOM METRIC -----------

		core.UpdateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Update info for the first custom metric and add new properties",
			},
			CustomMetricKey: common.FirstCustomMetric,
			CharacterKey:    common.CurrentCharacter,
			Case:            common.CreateProperties,
		},

		common.SaveToContextStage{
			Key:      common.FirstCustomMetric,
			JsonPath: "data.updateCustomMetric",
		},

		core.UpdateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Update existing properties and add one more",
			},
			CustomMetricKey: common.FirstCustomMetric,
			CharacterKey:    common.CurrentCharacter,
			Case:            common.UpdateProperties,
		},

		common.SaveToContextStage{
			Key:      common.FirstCustomMetric,
			JsonPath: "data.updateCustomMetric",
		},

		core.UpdateCustomMetricStage{
			Metadata: common.Metadata{
				Describe: "Delete the second property",
			},
			CustomMetricKey: common.FirstCustomMetric,
			CharacterKey:    common.CurrentCharacter,
			Case:            common.DeleteProperties,
		},

		common.SaveToContextStage{
			Key:      common.FirstCustomMetric,
			JsonPath: "data.updateCustomMetric",
		},

		core.ResetCustomMetricStage{
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

		// core.CreateMetricPropertyStage{
		// 	Metadata: common.Metadata{
		// 		Describe: "Create the first property for the second metric",
		// 	},
		// 	CharacterKey:    common.CurrentCharacter,
		// 	CustomMetricKey: common.SecondCustomMetric,
		// },

		// core.CreateMetricPropertyStage{
		// 	Metadata: common.Metadata{
		// 		Describe: "Create the second property for the second metric",
		// 	},
		// 	CharacterKey:    common.CurrentCharacter,
		// 	CustomMetricKey: common.SecondCustomMetric,
		// },

		// common.SaveToContextStage{
		// 	Key:      common.MetricProperty,
		// 	JsonPath: "data.createMetricProperty",
		// },

		// core.CreateMetricPropertyStage{
		// 	Metadata: common.Metadata{
		// 		Describe:    "Create the third property for the second metric",
		// 		ExpectError: true,
		// 	},
		// 	CharacterKey:    common.CurrentCharacter,
		// 	CustomMetricKey: common.SecondCustomMetric,
		// },

		// core.UpdateMetricPropertyStage{
		// 	Metadata: common.Metadata{
		// 		Describe: "Update the second property of the second metric",
		// 	},
		// 	CharacterKey:      common.CurrentCharacter,
		// 	CustomMetricKey:   common.SecondCustomMetric,
		// 	MetricPropertyKey: common.MetricProperty,
		// },

		// core.DeleteMetricPropertyStage{
		// 	Metadata: common.Metadata{
		// 		Describe: "Delete the second property of the second metric",
		// 	},
		// 	CharacterKey:      common.CurrentCharacter,
		// 	CustomMetricKey:   common.SecondCustomMetric,
		// 	MetricPropertyKey: common.MetricProperty,
		// },

		// core.DeleteCustomMetricStage{
		// 	Metadata: common.Metadata{
		// 		Describe: "Delete the second custom metric",
		// 	},
		// 	CustomMetricKey: common.SecondCustomMetric,
		// 	CharacterKey:    common.CurrentCharacter,
		// },

		core.GetProfileStage{
			Metadata: common.Metadata{
				Describe: "Get updated profile info after add characters, metrics, etc..",
			},
		},

		common.SaveToContextStage{
			Key:      common.Profile,
			JsonPath: "data.profile",
		},

		// --------- SNAPSHOT -----------

		analytics.CreateSnapshotStage{
			Metadata: common.Metadata{
				Describe: "Create a snapshot for the current character",
			},
			CharacterKey: common.CurrentCharacter,
		},

		common.SaveToContextStage{
			Key:      common.Snapshot,
			JsonPath: "data.createSnapshot",
		},

		analytics.GetSnapshotsStage{
			Metadata: common.Metadata{
				Describe: "Get all snapshots of the current character after add a new snapshot for it",
			},
			SnapshotKey:    common.Snapshot,
			HasOneSnapshot: true,
		},

		analytics.CreateSnapshotStage{
			Metadata: common.Metadata{
				Describe: "Create a snapshot for another character",
			},
			CharacterKey: common.AnotherCharacter,
		},

		analytics.GetSnapshotsStage{
			Metadata: common.Metadata{
				Describe: "Get all snapshots of the user after add a new snapshot for each of characters, so we have 2 snapshots",
			},
			HasTwoSnapshots: true,
		},

		// core.DeleteCharacterStage{
		// 	Metadata: common.Metadata{
		// 		Describe: "Delete another character",
		// 	},
		// 	CharacterKey: common.AnotherCharacter,
		// },

		core.GetProfileStage{
			Metadata: common.Metadata{
				Describe: "Get updated profile info for reviewing after performing flow",
			},
		},
	)

	err := p(&ctx)
	if err != nil {
		common.LogResponse()
	}

	if err == common.ErrEndStageReached {
		return
	}

	assert.Empty(t, err)
}
