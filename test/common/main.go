package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/steinfletcher/apitest"
	"github.com/tidwall/gjson"
)

type Metadata struct {
	Describe    string
	ExpectError bool
}

var jsonResponse string

func LogResponse() error {
	jsonData, err := json.MarshalIndent(gjson.Parse(jsonResponse).Value(), "", "  ")
	if err != nil {
		return err
	}

	log.Printf("--> Response: %v\n", string(jsonData))

	return nil
}

func ReadResponseJson(res *http.Response) string {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("failed to read response json")
	}
	defer res.Body.Close()

	return string(body)
}

type QueryParams struct {
	Url       string
	Query     string
	Variables map[string]interface{}
	Assertion func(*http.Response, *http.Request) error
}

func QueryGraphQL(ctx *context.Context, q *QueryParams) error {
	testingT, ok := (*ctx).Value(TestingT).(apitest.TestingT)
	if !ok {
		return ErrNotFoundInContext("TestingT")
	}

	if q.Url == "" {
		q.Url = GatewayUrl
	}

	result := apitest.New().
		EnableNetworking(cli).
		Post(q.Url).
		Header("Authorization", "Bearer "+IdToken).
		GraphQLQuery(q.Query, q.Variables).
		Expect(testingT).
		Status(http.StatusOK).
		Assert(q.Assertion).
		End()

	jsonResponse = ReadResponseJson(result.Response)
	LogResponse()

	return nil
}

type SaveToContextStage struct {
	Key      ContextKey
	JsonPath string
}

func (s SaveToContextStage) Exec(ctx *context.Context) error {
	value := gjson.Get(jsonResponse, s.JsonPath)
	if !value.Exists() {
		return fmt.Errorf("failed to get value with path %s", s.JsonPath)
	}

	*ctx = context.WithValue(*ctx, s.Key, value.Raw)
	return nil
}

// The stage to check if end stage is reached or not
type CheckEndStage struct {
	CurrentStage EndStage
}

func (s CheckEndStage) Exec(ctx *context.Context) error {
	// Get the end stage from the context
	endStage, ok := (*ctx).Value(EndStageKey).(string)
	if !ok {
		return ErrNotFoundInContext("EndStageKey")
	}

	// Check if the end stage is reached
	if s.CurrentStage == EndStage(endStage) {
		return ErrEndStageReached
	}

	return nil
}
