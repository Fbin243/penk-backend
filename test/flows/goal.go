package flows

import (
	"tenkhours/pkg/pineline"
	"tenkhours/test/common"
	"tenkhours/test/core"
)

var GoalFlow = []pineline.Stage{
	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Create the character for goal flow",
		},
		Case: common.CreateCharacter,
	},

	common.SaveToContextStage{
		Key:      common.Character,
		JsonPath: "data.upsertCharacter",
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Create the categories for goal flow",
		},
		Case:         common.CreateCategories,
		CharacterKey: common.Character,
	},

	common.SaveToContextStage{
		Key:      common.Character,
		JsonPath: "data.upsertCharacter",
	},

	core.UpsertCharacterStage{
		Metadata: common.Metadata{
			Describe: "Create the metrics for goal flow",
		},
		Case:         common.CreateMetrics,
		CharacterKey: common.Character,
	},

	common.SaveToContextStage{
		Key:      common.Character,
		JsonPath: "data.upsertCharacter",
	},

	core.UpsertGoal{
		Metadata: common.Metadata{
			Describe: "Create the goal",
		},
		Case:         common.CreateGoal,
		CharacterKey: common.Character,
	},

	common.SaveToContextStage{
		Key:      common.Goal,
		JsonPath: "data.upsertGoal",
	},

	core.UpsertGoal{
		Metadata: common.Metadata{
			Describe: "Update the goal",
		},
		Case:         common.UpdateGoal,
		GoalKey:      common.Goal,
		CharacterKey: common.Character,
	},

	core.DeleteGoal{
		Metadata: common.Metadata{
			Describe: "Delete the goal",
		},
		GoalKey: common.Goal,
	},
}
