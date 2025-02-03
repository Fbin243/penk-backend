package flows

import (
	"tenkhours/pkg/pineline"
	"tenkhours/test/common"
	"tenkhours/test/core"
)

var CharacterFlow = []pineline.Stage{
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

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Update the character",
		},
		CharacterKey: common.Character,
		Case:         common.UpdateCharacter,
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Create metrics in character",
		},
		CharacterKey:    common.Character,
		Case:            common.CreateMetrics,
		NumberOfMetrics: 5,
	},

	common.SaveToContextStage{
		Key:      common.Character,
		JsonPath: "data.upsertCharacter",
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Update metrics in character",
		},
		CharacterKey:    common.Character,
		Case:            common.UpdateMetrics,
		NumberOfMetrics: 5,
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Delete metrics in character",
		},
		CharacterKey:    common.Character,
		Case:            common.DeleteMetrics,
		NumberOfMetrics: 5,
	},

	core.DeleteCharacterStage{
		Metadata: common.Metadata{
			Describe: "Delete the character",
		},
		CharacterKey: common.Character,
	},
}
