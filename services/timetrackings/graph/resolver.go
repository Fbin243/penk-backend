package graph

import (
	"tenkhours/services/timetrackings/business"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TimeTrackingsBusiness *business.TimeTrackingsBusiness
}
