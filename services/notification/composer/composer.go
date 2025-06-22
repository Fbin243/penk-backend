package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/notification/business"
	mongorepo "tenkhours/services/notification/repo/mongo"

	"tenkhours/pkg/auth"
)

type Composer struct {
	DeviceTokenRepo business.IDeviceTokenRepo
	NotificationBiz business.INotificationBusiness
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	// Database
	db := mongodb.GetDBManager().DB

	// Repository
	devicesTokenRepo := mongorepo.NewDevicesTokenRepo(db)

	// Business
	firebaseManager := auth.GetFirebaseManager()

	notiBiz := business.NewNotificationBusiness(firebaseManager.MessagingClient, devicesTokenRepo)

	return &Composer{
		DeviceTokenRepo: devicesTokenRepo,
		NotificationBiz: notiBiz,
	}
}
