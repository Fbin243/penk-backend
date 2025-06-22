package composer

import (
	"tenkhours/services/notification/transport/graph"
)

func ComposeGraphQLResolver() *graph.Resolver {
	return &graph.Resolver{
		NotificationBusiness: GetComposer().NotificationBiz,
	}
}
