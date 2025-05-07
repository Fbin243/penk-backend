package messagequeue

import (
	"context"
	"log"

	"tenkhours/pkg/kafka"
	"tenkhours/services/notification/business"
	"tenkhours/services/notification/entity"
)

type NotificationQueue struct {
	producer         *kafka.Producer
	consumer         *kafka.Consumer
	notiBiz          *business.NotificationBusiness
	devicesTokenRepo business.IDeviceTokenRepo
}

func NewNotificationQueue(brokers []string, topic string, notiBiz *business.NotificationBusiness, devicesTokenRepo business.IDeviceTokenRepo) (*NotificationQueue, error) {
	// Create Kafka configuration
	cfg := kafka.DefaultConfig()
	cfg.Brokers = brokers

	// Create producer
	producer, err := kafka.NewProducer(cfg, topic)
	if err != nil {
		return nil, err
	}

	// Create consumer
	consumer, err := kafka.NewConsumer(cfg, topic, func(ctx context.Context, message []byte) error {
		var reminder entity.Reminder
		if err := kafka.UnmarshalMessage(message, &reminder); err != nil {
			return err
		}
		return processReminder(ctx, &reminder, notiBiz, devicesTokenRepo)
	})
	if err != nil {
		producer.Close()
		return nil, err
	}

	return &NotificationQueue{
		producer:         producer,
		consumer:         consumer,
		notiBiz:          notiBiz,
		devicesTokenRepo: devicesTokenRepo,
	}, nil
}

func (q *NotificationQueue) PublishReminder(ctx context.Context, reminder *entity.Reminder) error {
	return q.producer.Publish(ctx, reminder)
}

func (q *NotificationQueue) Start(ctx context.Context) error {
	return q.consumer.Start(ctx)
}

func (q *NotificationQueue) Close() error {
	if err := q.producer.Close(); err != nil {
		log.Printf("Error closing producer: %v", err)
	}
	if err := q.consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
	return nil
}

func processReminder(ctx context.Context, reminder *entity.Reminder, notiBiz *business.NotificationBusiness, devicesTokenRepo business.IDeviceTokenRepo) error {
	log.Printf("Processing Reminder: ID=%s, Title=%s, Time=%s", reminder.ID, reminder.Title, reminder.RemindTime)

	deviceIDs, err := devicesTokenRepo.GetDeviceIDsByProfileID(ctx, reminder.ProfileID)
	if err != nil {
		return err
	}

	for _, deviceID := range deviceIDs {
		req := &entity.SendNotiReq{
			ProfileID: reminder.ProfileID,
			DeviceID:  deviceID,
			Title:     "Reminder: " + reminder.Title,
			Body:      "It's time for your reminder!",
		}

		_, err = notiBiz.SendPushNotification(ctx, req)
		if err != nil {
			log.Printf("Failed to send notification to DeviceID=%s: %v", deviceID, err)
		} else {
			log.Printf("Notification sent successfully to DeviceID=%s", deviceID)
		}
	}

	return nil
}
