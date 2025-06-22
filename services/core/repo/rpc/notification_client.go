package rpc

import (
	"context"
	"log"

	"tenkhours/proto/pb/notification"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

type NotificationClient struct {
	notiClient notification.NotificationClient
}

func NewNotificationClient(notiClient notification.NotificationClient) *NotificationClient {
	return &NotificationClient{notiClient: notiClient}
}

func (c *NotificationClient) SendNotification(ctx context.Context, req *entity.SendNotiReq) (bool, error) {
	log.Print("Send request to Notification to send notification ...")
	notiReq := &notification.SendPushNotificationReq{}
	err := copier.Copy(notiReq, req)
	if err != nil {
		return false, err
	}

	_, err = c.notiClient.SendPushNotification(ctx, notiReq)
	if err != nil {
		return false, err
	}

	return true, nil
}
