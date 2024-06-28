package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func registerNewUser(t *testing.T, ctx TestContext) {
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(`mutation { registerAccount }`).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.data.registerAccount", true)).
		End()
}

func createNewCharacter(t *testing.T) {
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(`mutation {
			createCharacter(name: "Gymmer")
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.data.createCharacter", true)).
		End()
}

func createCustomMetrics(t *testing.T, ctx TestContext) {
	// Create the first custom metric with just name -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createCustomMetric(name: "calo", characterID: "%s") 
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.data.createCustomMetric", true)).
		End()

	// Create another custom metric with full information -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createCustomMetric(name: "workout", characterID: "%s", description: "workout description", style: {color: "red", icon: "i.png"})
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.data.createCustomMetric", true)).
		End()

	// Create the third custom metric to reach the limit -> failed
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createCustomMetric(name: "nutrition", characterID: "%s")
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.data.createCustomMetric", false)).
		End()
}

func createProperties(t *testing.T, ctx TestContext) {
	// Create the first property -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			updateCustomMetric(id: "%s", characterID: "%s", properties: [
				{name: "water", type: "Number", value: "20", unit: "l"}
			])
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.data.updateCustomMetric", true)).
		End()

	// Create the second --> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			updateCustomMetric(id: "%s", characterID: "%s", properties: [
				{name: "water", type: "Number", value: "20", unit: "l"},
				{name: "vitamin", type: "Number", value: "21", unit: "type"}
			])
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.data.updateCustomMetric", true)).
		End()

	// Create the third --> failed
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			updateCustomMetric(id: "%s", characterID: "%s", properties: [
				{name: "water", type: "Number", value: "20", unit: "l"},
				{name: "vitamin", type: "Number", value: "21", unit: "type"},
				{name: "vitamin", type: "Number", value: "21", unit: "type"}
			])
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.data.updateCustomMetric", false)).
		End()
}

func startTimeTracking(t *testing.T, ctx TestContext, trackWithMetric bool) {
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
		Assert(jsonpath.Equal("$.data.createTimeTracking", true)).
		End()
}

func stopTimeTracking(t *testing.T, ctx TestContext) {
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
		Assert(jsonpath.Equal("$.data.updateTimeTracking", true)).
		End()
}
