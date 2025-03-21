package composer

import (
	"tenkhours/services/core/transport/graph"
)

func ComposeGraphQLResolver() *graph.Resolver {
	composer := GetComposer()
	return &graph.Resolver{
		ProfileBusiness:   composer.ProfileBiz,
		CharacterBusiness: composer.CharacaterBiz,
		GoalBusiness:      composer.GoalBiz,
		CharacterRepo:     composer.CharacterRepo,
		MetricRepo:        composer.MetricRepo,
		CategoryRepo:      composer.CategoryRepo,
	}
}
