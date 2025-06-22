package rpc

import (
	"context"

	noti "tenkhours/proto/pb/notification"
	"tenkhours/services/notification/business"
	"tenkhours/services/notification/entity"

	"github.com/jinzhu/copier"
)

type NotificationHandler struct {
	noti.UnimplementedNotificationServer
	notificationBusiness business.INotificationBusiness
}

func NewNotificationHandler(notificationBiz business.INotificationBusiness) *NotificationHandler {
	return &NotificationHandler{
		notificationBusiness: notificationBiz,
	}
}

func (h *NotificationHandler) SendPushNotification(ctx context.Context, req *noti.SendPushNotificationReq) (*noti.SendPushNotificationResp, error) {
	res := &noti.SendPushNotificationResp{Success: false}

	notiReq := &entity.SendNotiReq{}
	err := copier.Copy(notiReq, req)
	if err != nil {
		return nil, err
	}

	_, err = h.notificationBusiness.SendPushNotification(ctx, notiReq)
	if err != nil {
		return res, err
	}

	res.Success = true
	return res, nil
}
