package composer

import (
	"log"
	"os"

	"tenkhours/pkg/auth"
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/kafka"
	"tenkhours/services/notification/business"
	message_queue "tenkhours/services/notification/message-queue"
	mongorepo "tenkhours/services/notification/repo/mongo"
	redisrepo "tenkhours/services/notification/repo/redis"
)

type Composer struct {
	DeviceTokenRepo business.IDeviceTokenRepo
	NotificationBiz business.INotificationBusiness
	kafkaProducer   *message_queue.KafkaProducer
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
	// TODO: Testing with localhost (need to change to kafka:9092 when running in docker)
	brokers := []string{"localhost:9092"}
	if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		brokers = []string{envBrokers}
	}

	// Create Kafka configuration
	cfg := kafka.DefaultConfig()
	cfg.Brokers = brokers

	firebaseManager := auth.GetFirebaseManager()
	notiProducer, err := message_queue.NewKafkaProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	notiBiz := business.NewNotificationBusiness(firebaseManager.MessagingClient, devicesTokenRepo, reminderRepo, reminderCache, notiProducer)

	composer = &Composer{
		DeviceTokenRepo: devicesTokenRepo,
		NotificationBiz: notiBiz,
		kafkaProducer:   notiProducer,
	}

	return composer
}

func (c *Composer) Close() error {
	return c.kafkaProducer.Close()
}
