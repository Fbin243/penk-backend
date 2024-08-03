package assertions

import (
	"tenkhours/test/variables"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var CreateCustomMetric = jsonpath.Chain().NotPresent("$.errors").
	Present("$.data.createCustomMetric.id").
	Equal("$.data.createCustomMetric.name", variables.CreateCustomMetric(nil)["name"]).
	Equal("$.data.createCustomMetric.description", variables.CreateCustomMetric(nil)["description"]).
	Equal("$.data.createCustomMetric.style", variables.CreateCustomMetric(nil)["style"]).
	Equal("$.data.createCustomMetric.time", float64(0)).
	Equal("$.data.createCustomMetric.limitedPropertyNumber", float64(2))
