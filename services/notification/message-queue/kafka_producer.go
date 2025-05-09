package message_queue

import (
	"context"
	"fmt"

	"tenkhours/pkg/kafka"
	"tenkhours/services/notification/entity"
)

// KafkaProducer implements INotificationProducer using Kafka
type KafkaProducer struct {
	producer *kafka.Producer
}

// NewKafkaProducer creates a new Kafka notification producer
func NewKafkaProducer(cfg *kafka.Config) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(cfg, "reminders")
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &KafkaProducer{
		producer: producer,
	}, nil
}

// SendNotification sends a notification message to Kafka
func (p *KafkaProducer) SendNotification(ctx context.Context, message *entity.NotificationMessage) error {
	return p.producer.Publish(ctx, message)
}

// Close closes the Kafka producer
func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
