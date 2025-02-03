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

	timetracking.CreateTimeTrackingStage{
		Metadata: common.Metadata{
			Describe: "Start a new tracking session of the current character",
		},
		CharacterKey: common.Character,
	},

	common.SaveToContextStage{
		Key:      common.TimeTracking,
		JsonPath: "data.createTimeTracking",
	},

	timetracking.CreateTimeTrackingStage{
		Metadata: common.Metadata{
			Describe:    "Start an existing session again",
			ExpectError: true,
		},
		CharacterKey: common.Character,
	},

	timetracking.UpdateTimeTracking{
		Metadata: common.Metadata{
			Describe: "End an existing tracking session of the current character",
		},
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
}
