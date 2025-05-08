package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// MessageHandler is a function type that processes messages from Kafka
type MessageHandler func(ctx context.Context, message []byte) error

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer sarama.Consumer
	topic    string
	handler  MessageHandler
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *Config, topic string, handler MessageHandler) (*Consumer, error) {
	saramaConfig, err := cfg.ToSaramaConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama config: %w", err)
	}

	consumer, err := sarama.NewConsumer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
		handler:  handler,
	}, nil
}

// Start begins consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return fmt.Errorf("failed to get partitions for topic %s: %w", c.topic, err)
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Failed to start consumer for partition %d: %v", partition, err)
			continue
		}

		go c.processPartition(ctx, pc)
	}

	log.Printf("Kafka consumer started for topic %s", c.topic)
	return nil
}

// processPartition processes messages from a single Kafka partition
func (c *Consumer) processPartition(ctx context.Context, pc sarama.PartitionConsumer) {
	for msg := range pc.Messages() {
		if err := c.handler(ctx, msg.Value); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	}
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	return c.consumer.Close()
}

// UnmarshalMessage is a helper function to unmarshal a message into a struct
func UnmarshalMessage(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
 