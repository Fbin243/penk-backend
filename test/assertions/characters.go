package assertions

import (
	"tenkhours/test/variables"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var CreateCharacter = func(userId interface{}) *jsonpath.AssertionChain {
	return jsonpath.Chain().NotPresent("$.errors").
		Present("$.data.createCharacter.id").
		Equal("$.data.createCharacter.name", variables.CreateCharacter["name"]).
		Equal("$.data.createCharacter.avatar", variables.CreateCharacter["avatar"]).
		Equal("$.data.createCharacter.gender", variables.CreateCharacter["gender"]).
		Equal("$.data.createCharacter.tags", variables.CreateCharacter["tags"]).
		Equal("$.data.createCharacter.limitedMetricNumber", float64(2)).
		Equal("$.data.createCharacter.totalFocusedTime", float64(0)).
		Equal("$.data.createCharacter.userID", userId)
}

var UpdateCharacter = func(characterId interface{}) *jsonpath.AssertionChain {
	return jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.updateCharacter.id", variables.UpdateCharacter(characterId)["id"]).
		Equal("$.data.updateCharacter.name", variables.UpdateCharacter(characterId)["name"]).
		Equal("$.data.updateCharacter.avatar", variables.UpdateCharacter(characterId)["avatar"]).
		Equal("$.data.updateCharacter.tags", variables.UpdateCharacter(characterId)["tags"])
}
