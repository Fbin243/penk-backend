package main

import (
	"log"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Unable to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Unable to create channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"hello", // Queue name
		false,   // Durable
		false,   // Auto-delete
		false,   // Exclusive
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		log.Fatalf("Unable to declare queue: %v", err)
	}

	// Send a message
	body := "Hello RabbitMQ!"
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key (queue name)
		false,  // Mandatory
		false,  // Immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Unable to send message: %v", err)
	}

	log.Println("Sent:", body)
}
