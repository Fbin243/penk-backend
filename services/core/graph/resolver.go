//go:generate go run github.com/99designs/gqlgen generate
package graph

import "tenkhours/services/core/business"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*business.ProfilesBusiness
	*business.CharactersBusiness
	*business.GoalsBusiness
}
