//go:generate go run github.com/99designs/gqlgen generate
package graph

import "tenkhours/pkg/core"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*core.ProfilesHandler
	*core.CharactersHandler
}
