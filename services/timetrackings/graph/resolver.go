package graph

import "tenkhours/pkg/timetrackings"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*timetrackings.TimeTrackingsHandler
}
