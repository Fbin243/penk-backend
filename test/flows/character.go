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
			Describe: "Create categories in character",
		},
		CharacterKey: common.Character,
		Case:         common.CreateCategories,
	},

	common.SaveToContextStage{
		Key:      common.Character,
		JsonPath: "data.upsertCharacter",
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Update categories in character",
		},
		CharacterKey: common.Character,
		Case:         common.UpdateCategories,
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Delete categories in character",
		},
		CharacterKey: common.Character,
		Case:         common.DeleteCategories,
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Create metrics in character",
		},
		CharacterKey: common.Character,
		Case:         common.CreateMetrics,
	},

	common.SaveToContextStage{
		Key:      common.Character,
		JsonPath: "data.upsertCharacter",
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Update metrics in character",
		},
		CharacterKey: common.Character,
		Case:         common.UpdateMetrics,
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Delete metrics in character",
		},
		CharacterKey: common.Character,
		Case:         common.DeleteMetrics,
	},

	core.DeleteCharacterStage{
		Metadata: common.Metadata{
			Describe: "Delete the character",
		},
		CharacterKey: common.Character,
	},
}
