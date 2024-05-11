package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQNotifier struct {
	Channel      *amqp.Channel
	ExchangeName string
}

func NewRabbitMQNotifier(channel *amqp.Channel, exchangeName string) *RabbitMQNotifier {
	return &RabbitMQNotifier{
		Channel:      channel,
		ExchangeName: exchangeName,
	}
}

func (n *RabbitMQNotifier) Notify(message []byte, routingKey string) error {
	var (
		contentType = "application/json"
		ctx         = context.Background()
	)
	err := n.Channel.PublishWithContext(
		ctx,
		n.ExchangeName, // Use the Exchange property of the struct
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        message,
		})

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
