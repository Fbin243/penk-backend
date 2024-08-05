package metrics

import (
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var CreateCustomMetricAssertion = jsonpath.Chain().NotPresent("$.errors").
	Present("$.data.createCustomMetric.id").
	Equal("$.data.createCustomMetric.name", CreateCustomMetricVariable(nil)["name"]).
	Equal("$.data.createCustomMetric.description", CreateCustomMetricVariable(nil)["description"]).
	Equal("$.data.createCustomMetric.style", CreateCustomMetricVariable(nil)["style"]).
	Equal("$.data.createCustomMetric.time", float64(0)).
	Equal("$.data.createCustomMetric.limitedPropertyNumber", float64(2))

var UpdateCustomMetricAssertion = jsonpath.Chain().NotPresent("$.errors").
	Equal("$.data.updateCustomMetric.name", UpdateCustomMetricVariable["name"]).
	Equal("$.data.updateCustomMetric.description", UpdateCustomMetricVariable["description"]).
	Equal("$.data.updateCustomMetric.style", UpdateCustomMetricVariable["style"])
