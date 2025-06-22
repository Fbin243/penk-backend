package flows

import (
	"tenkhours/pkg/pineline"
	"tenkhours/test/common"
	"tenkhours/test/core"
	"tenkhours/test/timetracking"
)

var TimeTrackingFlow = []pineline.Stage{
	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Create the character",
		},
		Case: common.CreateCharacter,
	},

	common.SaveToContextStage{
		Key:      common.Character,
		JsonPath: "data.upsertCharacter",
	},

	timetracking.GetCurrentTimeTrackingStage{
		Metadata: common.Metadata{
			Describe: "Get the current tracking session before creating a new one",
		},
		Case:         -common.CurrentTimeTrackingExist,
		CharacterKey: common.Character,
	},

	timetracking.CreateTimeTrackingStage{
		Metadata: common.Metadata{
			Describe: "Start a new tracking session without category",
		},
		Case:         common.TimeTrackingWithoutCategory,
		CharacterKey: common.Character,
	},

	common.SaveToContextStage{
		Key:      common.TimeTracking,
		JsonPath: "data.createTimeTracking",
	},

	timetracking.CreateTimeTrackingStage{
		Metadata: common.Metadata{
			Describe:    "Start a session again",
			ExpectError: true,
		},
		Case:         common.TimeTrackingWithoutCategory,
		CharacterKey: common.Character,
	},

	timetracking.GetCurrentTimeTrackingStage{
		Metadata: common.Metadata{
			Describe: "Get the current tracking session without category",
		},
		Case:         common.CurrentTimeTrackingExist,
		CharacterKey: common.Character,
	},

	timetracking.UpdateTimeTracking{
		Metadata: common.Metadata{
			Describe: "End an existing tracking session without category",
		},
		Case:            common.TimeTrackingWithoutCategory,
		TimeTrackingKey: common.TimeTracking,
		CharacterKey:    common.Character,
	},

	timetracking.UpdateTimeTracking{
		Metadata: common.Metadata{
			Describe:    "End an finished session again",
			ExpectError: true,
		},
		TimeTrackingKey: common.TimeTracking,
		CharacterKey:    common.Character,
	},

	timetracking.CreateTimeTrackingStage{
		Metadata: common.Metadata{
			Describe: "Start a new tracking session with category",
		},
		Case:         common.TimeTrackingWithCategory,
		CharacterKey: common.Character,
	},

	common.SaveToContextStage{
		Key:      common.TimeTracking,
		JsonPath: "data.createTimeTracking",
	},

	timetracking.GetCurrentTimeTrackingStage{
		Metadata: common.Metadata{
			Describe: "Get the current tracking session with category",
		},
		Case:         common.CurrentTimeTrackingExist,
		CharacterKey: common.Character,
	},

	timetracking.UpdateTimeTracking{
		Metadata: common.Metadata{
			Describe: "End an existing tracking session with category",
		},
		Case:            common.TimeTrackingWithCategory,
		TimeTrackingKey: common.TimeTracking,
		CharacterKey:    common.Character,
	},
}
