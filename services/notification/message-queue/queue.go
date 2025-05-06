package messagequeue

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	// Create a new Kafka producer configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // Ensure we get confirmation of successful sends
	config.Producer.Retry.Max = 5           // Retry up to 5 times if sending fails

	// Initialize the Kafka producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	// Return the KafkaProducer instance
	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// Publish sends a message to the Kafka
func (p *KafkaProducer) Publish(ctx context.Context, message interface{}) error {
	// Serialize the message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Create a Kafka producer message
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(body),
	}

	// Send the message to Kafka
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	// Log the successful send
	log.Printf("Message sent to Kafka: topic=%s, partition=%d, offset=%d", p.topic, partition, offset)
	return nil
}

// Close closes the Kafka producer.
func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
