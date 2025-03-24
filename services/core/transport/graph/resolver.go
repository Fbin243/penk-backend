//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"tenkhours/services/core/business"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ProfileBusiness   business.IProfileBusiness
	CharacterBusiness business.ICharacterBusiness
	GoalBusiness      business.IGoalBusiness
	MetricBusiness    business.IMetricBusiness
	CategoryBusiness  business.ICategoryBusiness
	CharacterRepo     business.ICharacterRepo
	MetricRepo        business.IMetricRepo
	CategoryRepo      business.ICategoryRepo
	TimeTrackingRepo  business.ITimeTrackingRepo
}
