package test

import (
	"testing"
	"time"
)

func TestUserFlow(t *testing.T) {
	// Step 0: Initialize the test context
	mTest := NewTestManager()
	mTest.InitContext()

	// Step 1: Register a new user
	// registerNewUser(t, ctx)

	// Step 2: Create a new character
	createNewCharacter(t)
	t.Log("--> Id Character: ", mTest.GetContext().IdCharacter)

	// Step 3: Create custom metrics
	createCustomMetrics(t, mTest.GetContext())
	t.Log("--> Id Custom Metric: ", mTest.GetContext().IdCustomMetric)

	// Step 4: Create properties
	createProperties(t, mTest.GetContext())
	// t.Log("--> Id Property: ", mTest.GetContext().IdProperty)

	// Step 5: Start a session without a custom metric
	startTimeTracking(t, mTest.GetContext(), false)
	t.Log("--> Id Time Tracking: ", mTest.GetContext().IdTimeTracking)

	// Simulate focus duration
	time.Sleep(5 * time.Second)

	// Step 6: End the session
	stopTimeTracking(t, mTest.GetContext())

	// Step 7: Start a session with a custom metric
	startTimeTracking(t, mTest.GetContext(), true)
	t.Log("--> Id Time Tracking: ", mTest.GetContext().IdTimeTracking)

	// Simulate focus duration
	time.Sleep(5 * time.Second)

	// Step 8: End the session
	stopTimeTracking(t, mTest.GetContext())
}
