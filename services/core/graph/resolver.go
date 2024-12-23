//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"tenkhours/services/core/business"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ProfilesBusiness   *business.ProfilesBusiness
	CharactersBusiness *business.CharactersBusiness
	GoalsBusiness      *business.GoalsBusiness
	TemplatesBusiness  *business.TemplatesBusiness
	SnapshotsBusiness  *business.SnapshotsBusiness
}
