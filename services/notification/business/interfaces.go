package business

import (
	"context"
	"time"

	"tenkhours/pkg/db/base"
	core_entity "tenkhours/services/core/entity"
	"tenkhours/services/notification/entity"
)

type INotificationBusiness interface {
	AddEmailToWaitlist(ctx context.Context, email string) error
	SendPushNotification(ctx context.Context, notiReq *entity.SendNotiReq) (bool, error)
	RegisterDeviceToken(ctx context.Context, profileID, token, deviceID, platform string) (bool, error)
	RemoveDeviceToken(ctx context.Context, profileID, token string) (bool, error)
	SyncTodayReminders(ctx context.Context) error
}

// repository
type IDeviceTokenRepo interface {
	base.IBaseRepo[entity.DevicesToken]
	UpsertDeviceToken(ctx context.Context, profileID, token, deviceID, platform string) error
	RemoveDeviceToken(ctx context.Context, profileID, token string) error
	GetDeviceTokenByDeviceID(ctx context.Context, deviceID string) (string, error)
	GetDeviceIDsByProfileID(ctx context.Context, profileID string) ([]string, error)
}

type IReminderRepo interface {
	base.IBaseRepo[core_entity.Reminder]
	GetTodayReminders(ctx context.Context) ([]core_entity.Reminder, error)
	BulkUpdateRemindTimes(ctx context.Context, reminders []core_entity.Reminder) error
	GetOutdatedReminders(ctx context.Context, now time.Time) ([]core_entity.Reminder, error)
}

type IReminderCache interface {
	GetReminder(ctx context.Context, id string) (*core_entity.Reminder, error)
	SetReminder(ctx context.Context, reminder *core_entity.Reminder) error
	DeleteReminder(ctx context.Context, id string) error
	SetReminders(ctx context.Context, reminders []core_entity.Reminder) error
	GetAllReminders(ctx context.Context) ([]core_entity.Reminder, error)
	ClearReminders(ctx context.Context) error
	GetRemindersWithMinScore(ctx context.Context) ([]core_entity.Reminder, error)
}
