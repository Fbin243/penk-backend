package test

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"tenkhours/pkg/pineline"
	"tenkhours/test/common"
	"tenkhours/test/core"
	"tenkhours/test/flows"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
)

func TestUserFlow(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.TestingT, t)

	token := os.Getenv("TOKEN")
	if token == "" {
		t.Fatalf("Token not set in environment variables")
	}

	log.Println("--> TOKEN: ", token)
	common.IdToken = token
	cleanUp := func(ctx *context.Context) {
		// Clean up the profile and related data
		err := common.QueryGraphQL(ctx,
			&common.QueryParams{
				Query:     core.DeleteProfileQuery,
				Assertion: []common.Assertion{jsonpath.Chain().NotPresent("$.errors").End()},
			})

		if err != nil {
			log.Fatalf("Failed to clean up the profile and related data %v\n", err)
		}
	}

	defer cleanUp(&ctx)

	// Get the flows to run the test
	flowKeyStr := os.Getenv("FLOWS")
	if flowKeyStr == "" {
		flowKeyStr = "profile"
	}

	flowsMap := map[common.FlowKey][]pineline.Stage{
		common.ProfilesFlowKey:      flows.ProfilesFlow,
		common.CharactersFlowKey:    flows.CharactersFlow,
		common.TimeTrackingsFlowKey: flows.TimeTrackingsFlow,
		common.AnalyticsFlowKey:     flows.AnalyticsFlow,
	}

	flowKeys := strings.Split(flowKeyStr, ",")
	log.Println("--> FLOWS: ", flowKeys)
	for _, flowKey := range flowKeys {
		p := pineline.Pineline(flowsMap[common.FlowKey(flowKey)]...)
		err := p(&ctx)

		if err != nil {
			common.LogResponse()
		}

		assert.Empty(t, err)
	}
}
