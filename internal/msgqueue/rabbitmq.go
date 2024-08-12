package msgqueue

import (
	"encoding/json"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"log"

	"github.com/streadway/amqp"
)

type rabbitMQ struct {
	exchange   string
	queue      string
	routingKey string
	rmqurl     string
	conn       *amqp.Connection
	channel    *amqp.Channel
}

func newRabbitMQ(config *models.RabbitMQConfig) models.MessageQueue {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Create a new channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// Declare a queue
	_, err = ch.QueueDeclare(
		config.Queue, // Queue name
		false,        // Durable
		false,        // Delete when unused
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Bind the queue to the exchange with a routing key
	routingKey := config.RoutingKey
	err = ch.QueueBind(
		config.Queue,    // Queue name
		routingKey,      // Routing key
		config.Exchange, // Exchange name
		false,           // No-wait
		nil,             // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind the queue to the exchange: %v", err)
	}

	return &rabbitMQ{queue: config.Queue, exchange: config.Exchange, routingKey: config.RoutingKey, rmqurl: config.URL, channel: ch, conn: conn}
}

func (r *rabbitMQ) Publish(event *models.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Publish a message to the exchange with a routing key
	err = r.channel.Publish(
		r.exchange,   // Exchange name
		r.routingKey, // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf("Published event: %s", event.Name)
	return nil
}
