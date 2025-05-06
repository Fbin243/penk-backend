package composer

import (
	"log"

	"tenkhours/pkg/auth"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/notification/business"
	"tenkhours/services/notification/cron"
	messagequeue "tenkhours/services/notification/message-queue"
	mongorepo "tenkhours/services/notification/repo/mongo"
)

type Composer struct {
	DeviceTokenRepo business.IDeviceTokenRepo
	NotificationBiz business.INotificationBusiness
	ReminderBiz     business.IReminderBusiness
	KafkaProducer   *messagequeue.KafkaProducer
	KafkaConsumer   *messagequeue.KafkaConsumer
	ReminderCron    *cron.ReminderCron
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	db := mongodb.GetDBManager().DB

	devicesTokenRepo := mongorepo.NewDevicesTokenRepo(db)
	reminderRepo := mongorepo.NewReminderRepo(db)

	kafkaProducer, err := messagequeue.NewKafkaProducer([]string{"localhost:9092"}, "reminder-topic")
	if err != nil {
		log.Fatalf("Failed to connect to Kafka producer: %v", err)
	}

	firebaseManager := auth.GetFirebaseManager()
	notiBiz := business.NewNotificationBusiness(firebaseManager.MessagingClient, devicesTokenRepo)
	reminderBiz := business.NewReminderBusiness(reminderRepo)

	kafkaConsumer, err := messagequeue.NewKafkaConsumer([]string{"localhost:9092"}, "reminder-topic", notiBiz, devicesTokenRepo)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka consumer: %v", err)
	}

	reminderCron := cron.NewReminderCron(reminderRepo, kafkaProducer)

	composer = &Composer{
		DeviceTokenRepo: devicesTokenRepo,
		NotificationBiz: notiBiz,
		ReminderBiz:     reminderBiz,
		KafkaProducer:   kafkaProducer,
		KafkaConsumer:   kafkaConsumer,
		ReminderCron:    reminderCron,
	}

	return composer
}
