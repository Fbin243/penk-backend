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
)

func TestUserFlow(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.TestingT, t)

	token := os.Getenv("TOKEN")
	if token == "" {
		t.Fatalf("Token not set in environment variables")
	}

	log.Println("--> TOKEN: ", token)
	common.IdToken = token
	common.DeviceId = "test-device-id"
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
		common.ProfileFlowKey:      flows.ProfileFlow,
		common.CharacterFlowKey:    flows.CharacterFlow,
		common.TimeTrackingFlowKey: flows.TimeTrackingFlow,
		common.SnapshotFlowKey:     flows.SnapshotFlow,
		common.GoalFlowKey:         flows.GoalFlow,
	}

	flowKeys := strings.Split(flowKeyStr, ",")
	log.Println("--> FLOWS: ", flowKeys)
	for _, flowKey := range flowKeys {
		err := pineline.Pineline(flowsMap[common.FlowKey(flowKey)]...)(&ctx)
		if err != nil {
			t.Fatalf("Failed to run the flow %v\n", flowKey)
			break
		}
	}
}
