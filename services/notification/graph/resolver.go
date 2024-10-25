package graph

import "tenkhours/pkg/notification"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	NotificationService notification.NotificationService
}
