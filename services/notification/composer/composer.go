package composer

import (
	"os"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/kafka"
	"tenkhours/services/notification/business"

	mongodb "tenkhours/pkg/db/mongo"

	mongorepo "tenkhours/services/notification/repo/mongo"
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

	db := mongodb.GetDBManager().DB

	devicesTokenRepo := mongorepo.NewDevicesTokenRepo(db)

	// Get Kafka brokers from environment or use default
	brokers := []string{"kafka:9092"}
	if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		brokers = []string{envBrokers}
	}

	// Create Kafka configuration
	cfg := kafka.DefaultConfig()
	cfg.Brokers = brokers

	firebaseManager := auth.GetFirebaseManager()
	notiBiz := business.NewNotificationBusiness(firebaseManager.MessagingClient, devicesTokenRepo)

	composer = &Composer{
		DeviceTokenRepo: devicesTokenRepo,
		NotificationBiz: notiBiz,
	}

	return composer
}
