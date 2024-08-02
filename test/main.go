package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"tenkhours/pineline"

	"github.com/steinfletcher/apitest"
	"github.com/yalp/jsonpath"
)

var response = map[string]interface{}{}

func logResponse(ctx *context.Context) error {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("--> Response: %v\n", string(jsonData))
	return nil
}

type QueryParams struct {
	Query          string
	Variables      map[string]interface{}
	AssertionChain func(*http.Response, *http.Request) error
}

func saveToContext(key ContextKey, jsonPath string) pineline.Stage {
	return func(ctx *context.Context) error {
		value, err := jsonpath.Read(response, jsonPath)
		if err != nil {
			return err
		}

		*ctx = context.WithValue(*ctx, key, value)
		return nil
	}
}

func queryGraphQL(queryParamsFunc func(ctx *context.Context) (*QueryParams, error)) pineline.Stage {
	return func(ctx *context.Context) error {
		testingT, ok := (*ctx).Value(TestingT).(apitest.TestingT)
		if !ok {
			return ErrNotFoundInContext(TestingT)
		}

		queryParams, err := queryParamsFunc(ctx)
		if err != nil {
			return fmt.Errorf("failed to make query params")
		}

		apitest.New().
			EnableNetworking(cli).
			Post(url).
			Header("Authorization", "Bearer "+IdToken).
			GraphQLQuery(queryParams.Query, queryParams.Variables).
			Expect(testingT).
			Status(http.StatusOK).
			Assert(queryParams.AssertionChain).
			End().
			JSON(&response)

		return nil
	}
}
