package test

import (
	"testing"
)

var ctx *TestContext = &TestContext{}

func TestCreate(t *testing.T) {
	registerNewUser(t, ctx)
	createNewCharacter(t, ctx)
	createCustomMetrics(t, ctx)
	createProperties(t, ctx)
}

func TestTimeTracking(t *testing.T) {
	t.Logf("Testing time tracking flow: %v", ctx)

	// Test time tracking without a custom metric
	startTimeTracking(t, ctx, false)
	startTimeTracking(t, ctx, false)
	stopTimeTracking(t, ctx)

	// Test time tracking with a custom metric
	startTimeTracking(t, ctx, true)
	stopTimeTracking(t, ctx)
}

func TestGetInfo(t *testing.T) {
	t.Logf("Testing read info flow: %v", ctx)

	getUserInfo(t, ctx)
	getUserCharacters(t, ctx)
}

func TestUpdate(t *testing.T) {
	t.Logf("Testing update flow: %v", ctx)

	updateCharacter(t, ctx)
	updateCustomMetric(t, ctx)
	updateProperty(t, ctx)
}

func TestDelete(t *testing.T) {
	t.Logf("Testing delete flow: %v", ctx)

	deleteProperty(t, ctx)
	deleteCustomMetric(t, ctx)
	deleteCharacter(t, ctx)
}
