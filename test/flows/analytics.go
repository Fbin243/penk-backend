package flows

import (
	"tenkhours/pkg/pineline"
	"tenkhours/test/analytics"
	"tenkhours/test/common"
)

var AnalyticsFlow = []pineline.Stage{
	analytics.CreateSnapshotStage{
		Metadata: common.Metadata{
			Describe: "Create a snapshot for the character",
		},
		CharacterKey: common.Character,
	},

	common.SaveToContextStage{
		Key:      common.Snapshot,
		JsonPath: "data.createSnapshot",
	},

	analytics.GetSnapshotsStage{
		Metadata: common.Metadata{
			Describe: "Get all snapshots of the character after add a new snapshot for it",
		},
		SnapshotKey:       common.Snapshot,
		NumberOfSnapshots: 1,
	},
}
