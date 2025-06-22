package composer

import "tenkhours/services/notification/transport/rpc"

func ComposeRPCHandler() *rpc.NotificationHandler {
	composer = GetComposer()
	return rpc.NewNotificationHandler(composer.NotificationBiz)
}
