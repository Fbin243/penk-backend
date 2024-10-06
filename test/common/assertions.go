package common

import jsonpath "github.com/steinfletcher/apitest-jsonpath"

var (
	AssertionSuccess = jsonpath.Chain().NotPresent("$.errors")
	AssertionError   = jsonpath.Chain().Present("$.errors")
)
