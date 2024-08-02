package test

import (
	"context"

	"tenkhours/pineline"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func createCustomMetric(expectError bool) pineline.Stage {
	return queryGraphQL(func(ctx *context.Context) (*QueryParams, error) {
		characterID, ok := (*ctx).Value(CharacterID).(string)
		if !ok {
			return nil, ErrNotFoundInContext(CharacterID)
		}

		variables := map[string]interface{}{
			"characterID": characterID,
			"name":        "Metric name",
			"description": "This is the custom metric description",
			"style": map[string]interface{}{
				"color": "#000000",
				"icon":  "icon.png",
			},
		}

		assertionChain := jsonpath.Chain().NotPresent("$.errors").
			Present("$.data.createCustomMetric.id").
			Equal("$.data.createCustomMetric.name", variables["name"]).
			Equal("$.data.createCustomMetric.description", variables["description"]).
			Equal("$.data.createCustomMetric.style", variables["style"]).
			Equal("$.data.createCustomMetric.time", float64(0)).
			Equal("$.data.createCustomMetric.limitedPropertyNumber", float64(2)).End()

		if expectError {
			assertionChain = jsonpath.Chain().Present("$.errors").End()
		}

		return &QueryParams{
			Query: `
			mutation CreateCustomMetric($characterID: String!, $name: String, $style: MetricStyleInput, $description: String) {
				createCustomMetric(
					characterID: $characterID
					input: { 
						name: $name, 
						style: $style, 
						description: $description
					}
				) {
					description
					id
					limitedPropertyNumber
					name
					time
					properties {
						id
						name
						type
						unit
						value
					}
					style {
						color
						icon
					}
				}
			}`,
			Variables:      variables,
			AssertionChain: assertionChain,
		}, nil
	})
}
