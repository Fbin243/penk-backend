package flows

import (
	"tenkhours/pkg/pineline"
	"tenkhours/test/common"
	"tenkhours/test/core"
)

var ProfileFlow = []pineline.Stage{
	core.UpsertProfileStage{
		Metadata: common.Metadata{
			Describe: "Create a new profile",
		},
		Case: common.GetProfile,
	},

	common.SaveToContextStage{
		Key:      common.Profile,
		JsonPath: "data.profile",
	},

	core.UpsertProfileStage{
		Metadata: common.Metadata{
			Describe: "Update profile info",
		},
		Case: common.UpdateProfile,
	},
}
