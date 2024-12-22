package composer

import (
	"tenkhours/services/notifications/business"
	"tenkhours/services/notifications/graph"
)

func ComposeGraphQLResolver() *graph.Resolver {
	notificationBiz := business.NewNotificationBusiness()

	return &graph.Resolver{
		NotificationBusiness: notificationBiz,
	}
}
