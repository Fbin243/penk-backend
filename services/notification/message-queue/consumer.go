package messagequeue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	biz "tenkhours/services/notification/business"
	"tenkhours/services/notification/entity"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer         sarama.Consumer
	topic            string
	NotiBiz          *biz.NotificationBusiness
	DevicesTokenRepo biz.IDeviceTokenRepo
}

// NewKafkaConsumer initializes a new Kafka consumer.
// - brokers: List of Kafka broker addresses (e.g., ["localhost:9092"])
// - topic: The Kafka topic to consume from (e.g., "reminder-topic")
// - notiBiz: Business logic for sending notifications
// - devicesTokenRepo: Repository for retrieving device tokens
func NewKafkaConsumer(brokers []string, topic string, notiBiz *biz.NotificationBusiness, devicesTokenRepo biz.IDeviceTokenRepo) (*KafkaConsumer, error) {
	// Create a new Kafka consumer configuration
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	// Return the KafkaConsumer instance
	return &KafkaConsumer{
		consumer:         consumer,
		topic:            topic,
		NotiBiz:          notiBiz,
		DevicesTokenRepo: devicesTokenRepo,
	}, nil
}

// Start begins listening to the Kafka topic and processing messages.
func (c *KafkaConsumer) Start() {
	// Get the list of partitions for the topic
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		log.Fatalf("Failed to get partitions for topic %s: %v", c.topic, err)
	}

	ctx := context.Background()

	// Start a consumer for each partition
	for _, partition := range partitions {
		// Create a partition consumer starting from the newest offset
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Failed to start consumer for partition %d: %v", partition, err)
			continue
		}

		// Process each partition in a separate goroutine
		go c.processPartition(ctx, pc)
	}

	log.Printf("Kafka consumer started for topic %s", c.topic)
}

// processPartition processes messages from a single Kafka partition.
func (c *KafkaConsumer) processPartition(ctx context.Context, pc sarama.PartitionConsumer) {
	// Loop over incoming messages from the partition
	for msg := range pc.Messages() {
		// Deserialize the message into a Reminder struct
		var reminder entity.Reminder
		if err := json.Unmarshal(msg.Value, &reminder); err != nil {
			log.Printf("Error decoding reminder: %v", err)
			continue
		}

		// Calculate the delay until the reminder's RemindTime
		delay := time.Until(reminder.RemindTime)
		if delay > 0 {
			// If the reminder is not due yet, wait until the RemindTime
			log.Printf("Waiting %v for reminder %s", delay, reminder.ID)
			time.Sleep(delay)
		}

		// Process the reminder (send notifications)
		c.processReminder(ctx, &reminder)
	}
}

// processReminder processes a reminder by sending push notifications to the user's devices.
func (c *KafkaConsumer) processReminder(ctx context.Context, reminder *entity.Reminder) {
	log.Printf("Processing Reminder: ID=%s, Title=%s, Time=%s", reminder.ID, reminder.Title, reminder.RemindTime)

	// Retrieve the device IDs for the user (ProfileID)
	deviceIDs, err := c.DevicesTokenRepo.GetDeviceIDsByProfileID(ctx, reminder.ProfileID)
	if err != nil {
		log.Printf("Failed to get device IDs for ProfileID=%s: %v", reminder.ProfileID, err)
		return
	}

	// Send a push notification to each device
	for _, deviceID := range deviceIDs {
		req := &entity.SendNotiReq{
			ProfileID: reminder.ProfileID,
			DeviceID:  deviceID,
			Title:     "Reminder: " + reminder.Title,
			Body:      "It's time for your reminder!",
		}

		// Send the push notification
		_, err = c.NotiBiz.SendPushNotification(ctx, req)
		if err != nil {
			log.Printf("Failed to send notification to DeviceID=%s: %v", deviceID, err)
		} else {
			log.Printf("Notification sent successfully to DeviceID=%s", deviceID)
		}
	}
}

// Close closes the Kafka consumer.
func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}
