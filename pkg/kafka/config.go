package kafka

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

// Config holds all Kafka configuration settings
type Config struct {
	// Broker configuration
	Brokers []string

	// Producer configuration
	Producer ProducerConfig

	// Consumer configuration
	Consumer ConsumerConfig

	// General configuration
	Version string
}

// ProducerConfig holds producer-specific configuration
type ProducerConfig struct {
	// RequiredAcks specifies the required number of acknowledgments from replicas
	RequiredAcks sarama.RequiredAcks

	// MaxRetries specifies the maximum number of retries for a message
	MaxRetries int

	// RetryBackoff specifies the backoff time between retries
	RetryBackoff time.Duration

	// ReturnSuccesses specifies whether to return success messages
	ReturnSuccesses bool

	// ReturnErrors specifies whether to return error messages
	ReturnErrors bool

	// Timeout specifies the maximum time to wait for a response
	Timeout time.Duration
}

// ConsumerConfig holds consumer-specific configuration
type ConsumerConfig struct {
	// GroupID specifies the consumer group ID
	GroupID string

	// InitialOffset specifies the initial offset to use if no offset was previously committed
	InitialOffset int64

	// MaxWaitTime specifies the maximum time to wait for a message
	MaxWaitTime time.Duration

	// MaxProcessingTime specifies the maximum time to process a message
	MaxProcessingTime time.Duration

	// ReturnErrors specifies whether to return error messages
	ReturnErrors bool
}

// DefaultConfig returns a default Kafka configuration
func DefaultConfig() *Config {
	return &Config{
		Brokers: []string{"localhost:9092"},
		Version: "2.8.0",
		Producer: ProducerConfig{
			RequiredAcks:    sarama.WaitForLocal,
			MaxRetries:      5,
			RetryBackoff:    100 * time.Millisecond,
			ReturnSuccesses: true,
			ReturnErrors:    true,
			Timeout:         5 * time.Second,
		},
		Consumer: ConsumerConfig{
			GroupID:           "default-group",
			InitialOffset:     sarama.OffsetNewest,
			MaxWaitTime:       250 * time.Millisecond,
			MaxProcessingTime: 100 * time.Millisecond,
			ReturnErrors:      true,
		},
	}
}

// LoadFromEnv loads Kafka configuration from environment variables
func LoadFromEnv() *Config {
	config := DefaultConfig()

	// Load brokers from environment
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		config.Brokers = []string{brokers}
	}

	// Load version from environment
	if version := os.Getenv("KAFKA_VERSION"); version != "" {
		config.Version = version
	}

	// Load producer configuration
	if retries := os.Getenv("KAFKA_PRODUCER_MAX_RETRIES"); retries != "" {
		if val, err := strconv.Atoi(retries); err == nil {
			config.Producer.MaxRetries = val
		}
	}

	if timeout := os.Getenv("KAFKA_PRODUCER_TIMEOUT"); timeout != "" {
		if val, err := time.ParseDuration(timeout); err == nil {
			config.Producer.Timeout = val
		}
	}

	// Load consumer configuration
	if groupID := os.Getenv("KAFKA_CONSUMER_GROUP_ID"); groupID != "" {
		config.Consumer.GroupID = groupID
	}

	if maxWait := os.Getenv("KAFKA_CONSUMER_MAX_WAIT_TIME"); maxWait != "" {
		if val, err := time.ParseDuration(maxWait); err == nil {
			config.Consumer.MaxWaitTime = val
		}
	}

	return config
}

// ToSaramaConfig converts the configuration to Sarama configuration
func (c *Config) ToSaramaConfig() (*sarama.Config, error) {
	config := sarama.NewConfig()

	// Set version
	version, err := sarama.ParseKafkaVersion(c.Version)
	if err != nil {
		return nil, fmt.Errorf("invalid kafka version: %w", err)
	}
	config.Version = version

	// Set producer configuration
	config.Producer.RequiredAcks = c.Producer.RequiredAcks
	config.Producer.Retry.Max = c.Producer.MaxRetries
	config.Producer.Retry.Backoff = c.Producer.RetryBackoff
	config.Producer.Return.Successes = c.Producer.ReturnSuccesses
	config.Producer.Return.Errors = c.Producer.ReturnErrors
	config.Producer.Timeout = c.Producer.Timeout

	// Set consumer configuration
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = c.Consumer.InitialOffset
	config.Consumer.MaxWaitTime = c.Consumer.MaxWaitTime
	config.Consumer.MaxProcessingTime = c.Consumer.MaxProcessingTime
	config.Consumer.Return.Errors = c.Consumer.ReturnErrors

	return config, nil
}
