package users

import (
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var NewUserAssertion = jsonpath.Chain().NotPresent("$.errors").
	Present("$.data.user.id").
	Present("$.data.user.firebaseUID").
	Equal("$.data.user.availableSnapshots", float64(2)).
	Equal("$.data.user.autoSnapshot", true).
	Equal("$.data.user.characters", []interface{}{}).
	Equal("$.data.user.currentCharacterID", nil)

var UpdateAccountAssertion = jsonpath.Chain().NotPresent("$.errors").
	Equal("$.data.updateAccount.name", UpdateAccountVariable["name"]).
	Equal("$.data.updateAccount.imageURL", UpdateAccountVariable["imageURL"]).
	Equal("$.data.updateAccount.currentCharacterID", UpdateAccountVariable["currentCharacterID"]).
	Equal("$.data.updateAccount.autoSnapshot", UpdateAccountVariable["autoSnapshot"])
