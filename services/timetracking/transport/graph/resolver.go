package graph

import (
	"tenkhours/services/timetracking/business"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TimeTrackingBusiness *business.TimeTrackingBusiness
}
