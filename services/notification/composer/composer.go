package composer

import (
	"log"
	"os"

	"tenkhours/pkg/auth"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/kafka"
	"tenkhours/services/notification/business"
	"tenkhours/services/notification/cron"
	messagequeue "tenkhours/services/notification/message-queue"
	mongorepo "tenkhours/services/notification/repo/mongo"
)

type Composer struct {
	DeviceTokenRepo   business.IDeviceTokenRepo
	NotificationBiz   business.INotificationBusiness
	ReminderBiz       business.IReminderBusiness
	NotificationQueue *messagequeue.NotificationQueue
	ReminderCron      *cron.ReminderCron
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	db := mongodb.GetDBManager().DB

	devicesTokenRepo := mongorepo.NewDevicesTokenRepo(db)
	reminderRepo := mongorepo.NewReminderRepo(db)

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
	reminderBiz := business.NewReminderBusiness(reminderRepo)

	// Create notification queue with both producer and consumer
	notificationQueue, err := messagequeue.NewNotificationQueue(brokers, "reminder-topic", notiBiz, devicesTokenRepo)
	if err != nil {
		log.Fatalf("Failed to create notification queue: %v", err)
	}

	reminderCron := cron.NewReminderCron(reminderRepo, notificationQueue)

	composer = &Composer{
		DeviceTokenRepo:   devicesTokenRepo,
		NotificationBiz:   notiBiz,
		ReminderBiz:       reminderBiz,
		NotificationQueue: notificationQueue,
		ReminderCron:      reminderCron,
	}

	return composer
}
