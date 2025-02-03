package composer

import (
	"tenkhours/services/timetracking/business"
	"tenkhours/services/timetracking/transport/graph"
)

func ComposeGraphQLResolver() *graph.Resolver {
	composer := GetComposer()
	timetrackingsBiz := business.NewTimeTrackingsBusiness(composer.coreClient, composer.currencyClient, composer.redisRepo)

	return &graph.Resolver{
		TimeTrackingBusiness: timetrackingsBiz,
	}
}
