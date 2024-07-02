package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func startTimeTracking(t *testing.T, ctx *TestContext, trackWithMetric bool) {
	gqlQuery := fmt.Sprintf(`mutation { 
		createTimeTracking(characterID: "%s")
	}`, ctx.IdCharacter)

	if trackWithMetric {
		gqlQuery = fmt.Sprintf(`mutation { 
			createTimeTracking(characterID: "%s", customMetricID: "%s")
		}`, ctx.IdCharacter, ctx.IdCustomMetric)
	}

	// Create a time tracking -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(gqlQuery).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		Assert(func(r1 *http.Response, r2 *http.Request) error {
			data, err := decodeResponseData(r1)
			if err != nil {
				return err
			}

			if idTimeTracking, ok := data["createTimeTracking"].(string); ok {
				ctx.IdTimeTracking = idTimeTracking
				return nil
			}

			return fmt.Errorf("failed to create a time tracking")
		}).
		Assert(logResponseData).
		End()
}

func stopTimeTracking(t *testing.T, ctx *TestContext) {
	// Stop the time tracking -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			updateTimeTracking(id: "%s")
		}`, ctx.IdTimeTracking)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		Assert(logResponseData).
		End()
}
