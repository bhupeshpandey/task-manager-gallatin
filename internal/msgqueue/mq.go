package msgqueue

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"log"
)

func NewMesageQueue(config *models.MessageQueueConfig) models.MessageQueue {
	var mq models.MessageQueue
	// Handle the message queue configuration based on the type and validate
	switch config.Type {
	case "RABBITMQ":
		if config.RabbitMQ != nil {
			fmt.Printf("Using rabbitMQ at %s\n", config.RabbitMQ.URL)
			fmt.Printf("exchange: %s, queue: %s, Routing Key: %s\n",
				config.RabbitMQ.Exchange,
				config.RabbitMQ.Queue,
				config.RabbitMQ.RoutingKey)
			mq = newRabbitMQ(config.RabbitMQ)
		} else {
			log.Fatalf("rabbitMQ configuration is missing")
		}
	case "KAFKA":
		// not handling right now.
		if config.Kafka != nil {
			fmt.Printf("Using Kafka with brokers: %v\n", config.Kafka.Brokers)
			fmt.Printf("Topic: %s, Group Id: %s\n", config.Kafka.Topic, config.Kafka.GroupID)
		} else {
			log.Fatalf("Kafka configuration is missing")
		}
	default:
		log.Fatalf("Unsupported message queue type: %s", config.Type)
	}

	return mq
}
