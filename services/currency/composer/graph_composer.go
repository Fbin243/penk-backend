package composer

import (
	"tenkhours/services/currency/transport/graph"
)

func ComposeGraphQLResolver() *graph.Resolver {
	composer := GetComposer()
	return &graph.Resolver{
		RewardBusiness: composer.RewardBiz,
	}
}
