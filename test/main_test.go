package test

import (
	"testing"
)

var ctx *TestContext = &TestContext{}

func TestCreate(t *testing.T) {
	// registerNewUser(t, ctx)
	createNewCharacter(t, ctx)
	createCustomMetrics(t, ctx)
	createProperties(t, ctx)
}

// func TestTimeTracking(t *testing.T) {
// 	// Test time tracking without a custom metric
// 	startTimeTracking(t, ctx, false)
// 	stopTimeTracking(t, ctx)

// 	// Test time tracking with a custom metric
// 	startTimeTracking(t, ctx, true)
// 	stopTimeTracking(t, ctx)
// }

func TestGetInfo(t *testing.T) {
	getUserInfo(t, ctx)
	getUserCharacters(t, ctx)
}

func TestUpdate(t *testing.T) {
	updateUser(t, ctx)
	updateCharacter(t, ctx)
	updateCustomMetric(t, ctx)
	updateProperty(t, ctx)
}

func TestDelete(t *testing.T) {
	deleteProperty(t, ctx)
	deleteCustomMetric(t, ctx)
	deleteCharacter(t, ctx)
}
