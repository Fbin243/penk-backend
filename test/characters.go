package test

import (
	"net/http"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func createNewCharacter(ctx *Context) error {
	testingT, ok := (*ctx)["testingT"].(apitest.TestingT)
	if !ok {
		return ErrNotFoundInContext
	}

	userID, ok := (*ctx)["userID"].(string)
	if !ok {
		return ErrNotFoundInContext
	}

	query := `
	mutation CreateCharacter($name: String!, $gender: Boolean, $avatar: String!, $tags: [String]) {
		createCharacter(
			input: { 
				name: $name, 
				gender: $gender, 
				avatar: $avatar, 
				tags: $tags 
			}
		) {
			avatar
			gender
			id
			limitedMetricNumber
			name
			tags
			totalFocusedTime
			userID
		}
	}`

	variables := Map{
		"name":   "Conan",
		"gender": true,
		"avatar": "avatar.png",
		"tags":   []interface{}{"#Kid", "#Detective"},
	}

	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(query, variables).
		Expect(testingT).
		Status(http.StatusOK).
		Assert(jsonpath.NotPresent("$.errors")).
		Assert(jsonpath.Equal("$.data.createCharacter.name", variables["name"])).
		Assert(jsonpath.Equal("$.data.createCharacter.avatar", variables["avatar"])).
		Assert(jsonpath.Equal("$.data.createCharacter.gender", variables["gender"])).
		Assert(jsonpath.Equal("$.data.createCharacter.tags", variables["tags"])).
		Assert(jsonpath.Present("$.data.createCharacter.id")).
		Assert(jsonpath.Equal("$.data.createCharacter.limitedMetricNumber", float64(2))).
		Assert(jsonpath.Equal("$.data.createCharacter.totalFocusedTime", float64(0))).
		Assert(jsonpath.Equal("$.data.createCharacter.userID", userID)).
		End().JSON(response)

	response.log()
	return nil
}

// func createCustomMetrics(t *testing.T, ctx *TestContext) {
// 	// Create the first custom metric with just name -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			createCustomMetric(name: "Test metric 1", characterID: "%s") {id}
// 		}`, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	ctx.IdCustomMetric = ResponseBody["data"].(map[string]interface{})["createCustomMetric"].(map[string]interface{})["id"].(string)
// 	logResponse(ResponseBody)

// 	// Create another custom metric with full information -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			createCustomMetric(name: "Test metric 2", characterID: "%s", description: "Test metric description", style: {color: "#000000", icon: "icon.png"}) {id}
// 		}`, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)

// 	// Create the third custom metric to reach the limit -> failed
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			createCustomMetric(name: "Test metric 3", characterID: "%s") {id}
// 		}`, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.Present("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)
// }

// func createProperties(t *testing.T, ctx *TestContext) {
// 	// Create the first property -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			createMetricProperty(metricID: "%s", characterID: "%s", name: "Test property 1", type: "Number", value: "10", unit: "kg") {id}
// 		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	ctx.IdProperty = ResponseBody["data"].(map[string]interface{})["createMetricProperty"].(map[string]interface{})["id"].(string)
// 	logResponse(ResponseBody)

// 	// Create the second --> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			createMetricProperty(metricID: "%s", characterID: "%s", name: "Test property 2", type: "Number", value: "20", unit: "l") {id}
// 		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)

// 	// Create the third --> failed
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			createMetricProperty(metricID: "%s", characterID: "%s", name: "Test property 3", type: "Number", value: "30", unit: "m") {id}
// 		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.Present("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)
// }

// /**
//  * READ
//  */
// func getUserCharacters(t *testing.T, ctx *TestContext) {
// 	// Get the user's character -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(`query {
// 			userCharacters { id limitedMetricNumber name gender tags totalFocusedTime userID customMetrics { description id limitedPropertyNumber name time properties { id name type unit value } style { color icon } } }
// 		}`).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)
// }

// /**
//  * UPDATE
//  */
// func updateCharacter(t *testing.T, ctx *TestContext) {
// 	// Update the character's name -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			updateCharacter(id: "%s", name: "Updated Character") {id}
// 		}`, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)
// }

// func updateCustomMetric(t *testing.T, ctx *TestContext) {
// 	// Update the custom metric's name -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			updateCustomMetric(id: "%s", characterID: "%s" ,name: "Updated Metric") {id}
// 		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)
// }

// func updateProperty(t *testing.T, ctx *TestContext) {
// 	// Update the property's value -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			updateMetricProperty(id: "%s", metricID: "%s", characterID: "%s", name: "Updated property", value: "50") {id}
// 		}`, ctx.IdProperty, ctx.IdCustomMetric, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)
// }

// func deleteProperty(t *testing.T, ctx *TestContext) {
// 	// Delete the property -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			deleteMetricProperty(id: "%s", metricID: "%s", characterID: "%s") {id}
// 		}`, ctx.IdProperty, ctx.IdCustomMetric, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)
// }

// func deleteCustomMetric(t *testing.T, ctx *TestContext) {
// 	// Delete the custom metric -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			deleteCustomMetric(id: "%s", characterID: "%s") {id}
// 		}`, ctx.IdCustomMetric, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)
// }

// func deleteCharacter(t *testing.T, ctx *TestContext) {
// 	// Delete the character -> success
// 	apitest.New().
// 		EnableNetworking(cli).
// 		Post(url).
// 		Header("Authorization", "Bearer "+IdToken).
// 		GraphQLQuery(fmt.Sprintf(`mutation {
// 			deleteCharacter(id: "%s") {id}
// 		}`, ctx.IdCharacter)).
// 		Expect(t).
// 		Status(http.StatusOK).
// 		Assert(jsonpath.NotPresent("$.errors")).
// 		End().JSON(&ResponseBody)

// 	logResponse(ResponseBody)
// }
