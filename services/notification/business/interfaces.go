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
	GetDeviceIDsByProfileID(ctx context.Context, profileID string) ([]string, error)
}

type IReminderBusiness interface {
	CreateReminder(ctx context.Context, reminder entity.ReminderInput) (*entity.Reminder, error)
	GetRemindersByProfileID(ctx context.Context) ([]*entity.Reminder, error)
	GetReminderByID(ctx context.Context, reminderID string) (*entity.Reminder, error)
	UpdateReminder(ctx context.Context, reminder entity.ReminderInput) (*entity.Reminder, error)
	DeleteReminder(ctx context.Context, reminderID string) (bool, error)
}

type IReminderRepo interface {
	base.IBaseRepo[entity.Reminder]
	CreateReminder(ctx context.Context, reminder *entity.Reminder) (*entity.Reminder, error)
	GetRemindersByProfileID(ctx context.Context, profileID string) ([]*entity.Reminder, error)
	GetReminderByID(ctx context.Context, reminderID string) (*entity.Reminder, error)
	UpdateReminder(ctx context.Context, reminder *entity.Reminder) (*entity.Reminder, error)
	DeleteReminder(ctx context.Context, reminderID string) (bool, error)
	GetUpcomingReminders(ctx context.Context) ([]*entity.Reminder, error)
}
