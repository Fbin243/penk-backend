package metrics

import (
	"context"
	"log"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/tidwall/gjson"
)

type CreateCustomMetricStage struct {
	common.Metadata
	CharacterKey common.ContextKey // The key to the character in the context
}

func (s CreateCustomMetricStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFound(s.CharacterKey)
	}

	variables := map[string]interface{}{
		"characterID": gjson.Get(character, "id").Value(),
		"name":        "Metric name",
		"description": "This is the custom metric description",
		"style": map[string]interface{}{
			"color": "#000000",
			"icon":  "icon.png",
		},
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Present("$.data.createCustomMetric.id").
		Equal("$.data.createCustomMetric.name", variables["name"]).
		Equal("$.data.createCustomMetric.description", variables["description"]).
		Equal("$.data.createCustomMetric.style", variables["style"]).
		Equal("$.data.createCustomMetric.time", float64(0)).
		Equal("$.data.createCustomMetric.limitedPropertyNumber", float64(2))

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.CoreUrl,
			Query:     CreateCustomMetricQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

// map[string]interface{}{
// 	"name":        "Update metric name",
// 	"description": "This is the updated custom metric description",
// 	"style": map[string]interface{}{
// 		"color": "#123456",
// 		"icon":  "update_icon.png",
// 	},
// }
