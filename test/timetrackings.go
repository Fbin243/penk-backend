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
		createTimeTracking(characterID: "%s", startTime: time.Now()) {id}
	}`, ctx.IdCharacter)

	if trackWithMetric {
		gqlQuery = fmt.Sprintf(`mutation { 
			createTimeTracking(characterID: "%s", customMetricID: "%s") {id}
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
		End().JSON(&responseBody)

	ctx.IdTimeTracking = responseBody["data"].(map[string]interface{})["createTimeTracking"].(map[string]interface{})["id"].(string)
	logResponse(responseBody)
}

func stopTimeTracking(t *testing.T, ctx *TestContext) {
	// Stop the time tracking -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			updateTimeTracking(id: "%s") {id}
		}`, ctx.IdTimeTracking)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}
