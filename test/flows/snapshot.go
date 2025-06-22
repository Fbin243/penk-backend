package flows

import (
	"tenkhours/pkg/pineline"
	"tenkhours/test/common"
	"tenkhours/test/core"
)

var SnapshotFlow = []pineline.Stage{
	core.CreateSnapshotStage{
		Metadata: common.Metadata{
			Describe: "Create a snapshot for the character",
		},
		CharacterKey: common.Character,
	},

	common.SaveToContextStage{
		Key:      common.Snapshot,
		JsonPath: "data.createSnapshot",
	},

	core.GetSnapshotsStage{
		Metadata: common.Metadata{
			Describe: "Get all snapshots of the character after add a new snapshot for it",
		},
		SnapshotKey:       common.Snapshot,
		NumberOfSnapshots: 1,
	},
}
