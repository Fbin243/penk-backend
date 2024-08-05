package characters

import (
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var CreateCharacterAssertion = func(userId interface{}) *jsonpath.AssertionChain {
	return jsonpath.Chain().NotPresent("$.errors").
		Present("$.data.createCharacter.id").
		Equal("$.data.createCharacter.name", CreateCharacterVariable["name"]).
		Equal("$.data.createCharacter.avatar", CreateCharacterVariable["avatar"]).
		Equal("$.data.createCharacter.gender", CreateCharacterVariable["gender"]).
		Equal("$.data.createCharacter.tags", CreateCharacterVariable["tags"]).
		Equal("$.data.createCharacter.limitedMetricNumber", float64(2)).
		Equal("$.data.createCharacter.totalFocusedTime", float64(0)).
		Equal("$.data.createCharacter.userID", userId).
		Equal("$.data.createCharacter.customMetrics", []interface{}{})
}

var UpdateCharacterAssertion = func(characterId interface{}) *jsonpath.AssertionChain {
	return jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.updateCharacter.id", UpdateCharacterVariable(characterId)["id"]).
		Equal("$.data.updateCharacter.name", UpdateCharacterVariable(characterId)["name"]).
		Equal("$.data.updateCharacter.avatar", UpdateCharacterVariable(characterId)["avatar"]).
		Equal("$.data.updateCharacter.tags", UpdateCharacterVariable(characterId)["tags"])
}
