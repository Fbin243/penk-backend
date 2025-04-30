package business

import (
	"context"

	"tenkhours/pkg/db/base"
	"tenkhours/services/notification/entity"
)

type INotificationBusiness interface {
	AddEmailToWaitlist(ctx context.Context, email string) error
	SendPushNotification(ctx context.Context, notiReq *entity.SendNotiReq) (bool, error)
	RegisterDeviceToken(ctx context.Context, profileID, token, deviceID, platform string) (bool, error)
	RemoveDeviceToken(ctx context.Context, profileID, token string) (bool, error)
}

// repository
type IDeviceTokenRepo interface {
	base.IBaseRepo[entity.DevicesToken]
	UpsertDeviceToken(ctx context.Context, profileID, token, deviceID, platform string) error
	RemoveDeviceToken(ctx context.Context, profileID, token string) error
	GetDeviceTokenByDeviceID(ctx context.Context, deviceID string) (string, error)
}
