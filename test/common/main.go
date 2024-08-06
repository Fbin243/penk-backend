package common

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/steinfletcher/apitest"
	"github.com/yalp/jsonpath"
)

type Metadata struct {
	Name        string
	ExpectError bool
}

var response = map[string]interface{}{}

type SwitchUrlStage struct {
	NewUrl string
}

func (s SwitchUrlStage) Exec(ctx *context.Context) error {
	url = s.NewUrl
	return nil
}

func LogResponse() error {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("--> Response: %v\n", string(jsonData))

	return nil
}

type SaveToContextStage struct {
	Key      ContextKey
	JsonPath string
}

func (s SaveToContextStage) Exec(ctx *context.Context) error {
	value, err := jsonpath.Read(response, s.JsonPath)
	if err != nil {
		return err
	}

	*ctx = context.WithValue(*ctx, s.Key, value)
	return nil
}

type QueryParams struct {
	Query     string
	Variables map[string]interface{}
	Assertion func(*http.Response, *http.Request) error
}

func QueryGraphQL(ctx *context.Context, q *QueryParams) error {
	testingT, ok := (*ctx).Value(TestingT).(apitest.TestingT)
	if !ok {
		return ErrNotFoundInContext(TestingT)
	}

	// Reset response
	response = map[string]interface{}{}

	apitest.New().
		EnableNetworking(cli).
		Post(url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(q.Query, q.Variables).
		Expect(testingT).
		Status(http.StatusOK).
		Assert(q.Assertion).
		End().
		JSON(&response)

	LogResponse()

	return nil
}
