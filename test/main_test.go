package test

import (
	"testing"
	"time"
)

var ctx *TestContext = &TestContext{}

func TestCreateFlow(t *testing.T) {
	// registerNewUser(t, ctx)
	createNewCharacter(t, ctx)
	createCustomMetrics(t, ctx)
	createProperties(t, ctx)
}

func TestTimeTrackingFlow(t *testing.T) {
	t.Logf("Testing time tracking flow: %v", ctx)

	// Test time tracking without a custom metric
	startTimeTracking(t, ctx, false)
	time.Sleep(5 * time.Second)
	stopTimeTracking(t, ctx)

	// Test time tracking with a custom metric
	startTimeTracking(t, ctx, true)
	time.Sleep(5 * time.Second)
	stopTimeTracking(t, ctx)
}

// func TestUpdateFlow(t *testing.T) {
// 	updateCharacter(t, ctx)
// 	updateCustomMetrics(t, ctx)
// 	updateProperties(t, ctx)
// }
