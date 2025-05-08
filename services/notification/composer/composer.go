package composer

import (
	"os"

	"tenkhours/pkg/auth"
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/kafka"
	"tenkhours/services/notification/business"
	mongorepo "tenkhours/services/notification/repo/mongo"
	redisrepo "tenkhours/services/notification/repo/redis"
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
	redisClient := rdb.GetRedisClient()

	devicesTokenRepo := mongorepo.NewDevicesTokenRepo(db)
	reminderRepo := mongorepo.NewReminderRepo(db)
	reminderCache := redisrepo.NewReminderCache(redisClient)

	// Get Kafka brokers from environment or use default
	brokers := []string{"kafka:9092"}
	if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		brokers = []string{envBrokers}
	}

	// Create Kafka configuration
	cfg := kafka.DefaultConfig()
	cfg.Brokers = brokers

	firebaseManager := auth.GetFirebaseManager()
	notiBiz := business.NewNotificationBusiness(firebaseManager.MessagingClient, devicesTokenRepo, reminderRepo, reminderCache)

	composer = &Composer{
		DeviceTokenRepo: devicesTokenRepo,
		NotificationBiz: notiBiz,
	}

	return composer
}
