package composer

import (
	"tenkhours/services/analytic/transport/graph"
)

func ComposeGraphQLResolver() *graph.Resolver {
	return &graph.Resolver{
		AnalyticBusiness: GetComposer().AnalyticBiz,
	}
}
