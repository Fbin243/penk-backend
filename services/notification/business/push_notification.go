package business

import (
	"context"
	"fmt"
	"log"

	"tenkhours/pkg/utils"
	"tenkhours/services/notification/entity"

	"firebase.google.com/go/messaging"
)

// SendPushNotification sends a message to a device token
func (biz *NotificationBusiness) SendPushNotification(ctx context.Context, req *entity.SendNotiReq) (bool, error) {
	log.Printf("Req: %v", utils.PrettyJSON(req))
	deviceToken, err := biz.DevicesTokenRepo.GetDeviceTokenByDeviceID(ctx, req.DeviceID)
	if err != nil {
		return false, fmt.Errorf("failed to get device token: %v", err)
	}

	message := &messaging.Message{
		Token: deviceToken,
		Notification: &messaging.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
	}

	response, err := biz.MessagingClient.Send(ctx, message)
	if err != nil {
		return false, fmt.Errorf("failed to send push notification: %v", err)
	}

	log.Printf("Push notification sent successfully: %s", response)
	return true, nil
}
