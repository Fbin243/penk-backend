package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func registerNewUser(t *testing.T, ctx *TestContext) {
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(`mutation { registerAccount }`).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		Assert(func(r1 *http.Response, r2 *http.Request) error {
			data, err := decodeResponseData(r1)
			if err != nil {
				return err
			}

			if idUser, ok := data["registerAccount"].(string); ok {
				ctx.IdUser = idUser
				return nil
			}

			return fmt.Errorf("failed to register a new user")
		}).
		End()
}

func createNewCharacter(t *testing.T, ctx *TestContext) {
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(`mutation {
			createCharacter(name: "Test Character")
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		Assert(func(r1 *http.Response, r2 *http.Request) error {
			data, err := decodeResponseData(r1)
			if err != nil {
				return err
			}

			if idCharacter, ok := data["createCharacter"].(string); ok {
				ctx.IdCharacter = idCharacter
				return nil
			}

			return fmt.Errorf("failed to create a new character")
		}).
		End()
}

func createCustomMetrics(t *testing.T, ctx *TestContext) {
	// Create the first custom metric with just name -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createCustomMetric(name: "Test metric 1", characterID: "%s") 
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		Assert(func(r1 *http.Response, r2 *http.Request) error {
			data, err := decodeResponseData(r1)
			if err != nil {
				return err
			}

			if idCustomMetric, ok := data["createCustomMetric"].(string); ok {
				ctx.IdCustomMetric = idCustomMetric
				return nil
			}

			return fmt.Errorf("failed to create a custom metric")
		}).
		End()

	// Create another custom metric with full information -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createCustomMetric(name: "Test metric 2", characterID: "%s", description: "Test metric description", style: {color: "red", icon: "icon.png"})
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End()

	// Create the third custom metric to reach the limit -> failed
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createCustomMetric(name: "Test metric 3", characterID: "%s")
		}`, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.errors")).
		End()
}

func createProperties(t *testing.T, ctx *TestContext) {
	// Create the first property -> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createMetricProperty(metricID: "%s", characterID: "%s", name: "Test property 1", type: "Number", value: "10", unit: "kg")
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		Assert(func(r1 *http.Response, r2 *http.Request) error {
			data, err := decodeResponseData(r1)
			if err != nil {
				return err
			}

			if idProperty, ok := data["createMetricProperty"].(string); ok {
				ctx.IdProperty = idProperty
				return nil
			}

			return fmt.Errorf("failed to create a property")
		}).
		End()

	// Create the second --> success
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createMetricProperty(metricID: "%s", characterID: "%s", name: "Test property 2", type: "Number", value: "20", unit: "l")
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		End()

	// Create the third --> failed
	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(fmt.Sprintf(`mutation { 
			createMetricProperty(metricID: "%s", characterID: "%s", name: "Test property 3", type: "Number", value: "30", unit: "m")
		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.errors")).
		End()
}

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
		End()
}
