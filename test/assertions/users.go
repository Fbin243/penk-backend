package assertions

import (
	"tenkhours/test/variables"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var (
	AssertionSuccess = jsonpath.Chain().NotPresent("$.errors")
	AssertionError   = jsonpath.Chain().Present("$.errors")
)

var UpdateAccount = jsonpath.Chain().NotPresent("$.errors").
	Equal("$.data.updateAccount.name", variables.UpdateAccount["name"]).
	Equal("$.data.updateAccount.imageURL", variables.UpdateAccount["imageURL"]).
	Equal("$.data.updateAccount.currentCharacterID", variables.UpdateAccount["currentCharacterID"]).
	Equal("$.data.updateAccount.autoSnapshot", variables.UpdateAccount["autoSnapshot"])
