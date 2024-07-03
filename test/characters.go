package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

/**
 * CREATE
 */
func createNewCharacter(t *testing.T, ctx *TestContext) {
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(`mutation {
			createCharacter(name: "Test Character", gender: "true") {id}
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	ctx.IdCharacter = responseBody["data"].(map[string]interface{})["createCharacter"].(map[string]interface{})["id"].(string)

	logResponse(responseBody)
}

func createCustomMetrics(t *testing.T, ctx *TestContext) {
	// Create the first custom metric with just name -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createCustomMetric(name: "Test metric 1", characterID: "%s") {id}
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	ctx.IdCustomMetric = responseBody["data"].(map[string]interface{})["createCustomMetric"].(map[string]interface{})["id"].(string)
	logResponse(responseBody)

	// Create another custom metric with full information -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createCustomMetric(name: "Test metric 2", characterID: "%s", description: "Test metric description", style: {color: "#000000", icon: "icon.png"}) {id}
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)

	// Create the third custom metric to reach the limit -> failed
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createCustomMetric(name: "Test metric 3", characterID: "%s") {id}
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}

func createProperties(t *testing.T, ctx *TestContext) {
	// Create the first property -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createMetricProperty(metricID: "%s", characterID: "%s", name: "Test property 1", type: "Number", value: "10", unit: "kg") {id}
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	ctx.IdProperty = responseBody["data"].(map[string]interface{})["createMetricProperty"].(map[string]interface{})["id"].(string)
	logResponse(responseBody)

	// Create the second --> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createMetricProperty(metricID: "%s", characterID: "%s", name: "Test property 2", type: "Number", value: "20", unit: "l") {id}
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)

	// Create the third --> failed
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createMetricProperty(metricID: "%s", characterID: "%s", name: "Test property 3", type: "Number", value: "30", unit: "m") {id}
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}

/**
 * READ
 */
func getUserCharacters(t *testing.T, ctx *TestContext) {
	// Get the user's character -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(`query { 
			userCharacters { id limitedMetricNumber name gender tags totalFocusedTime userID customMetrics { description id limitedPropertyNumber name time properties { id name type unit value } style { color icon } } }
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}

/**
 * UPDATE
 */
func updateCharacter(t *testing.T, ctx *TestContext) {
	// Update the character's name -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			updateCharacter(id: "%s", name: "Updated Character") {id}
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}

func updateCustomMetric(t *testing.T, ctx *TestContext) {
	// Update the custom metric's name -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			updateCustomMetric(id: "%s", characterID: "%s" ,name: "Updated Metric") {id}
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}

func updateProperty(t *testing.T, ctx *TestContext) {
	// Update the property's value -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			updateMetricProperty(id: "%s", metricID: "%s", characterID: "%s", name: "Updated property", value: "50") {id}
		}`, ctx.IdProperty, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}

func deleteProperty(t *testing.T, ctx *TestContext) {
	// Delete the property -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			deleteMetricProperty(id: "%s", metricID: "%s", characterID: "%s") {id}
		}`, ctx.IdProperty, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}

func deleteCustomMetric(t *testing.T, ctx *TestContext) {
	// Delete the custom metric -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			deleteCustomMetric(id: "%s", characterID: "%s") {id}
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}

func deleteCharacter(t *testing.T, ctx *TestContext) {
	// Delete the character -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			deleteCharacter(id: "%s") {id}
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End().JSON(&responseBody)

	logResponse(responseBody)
}
