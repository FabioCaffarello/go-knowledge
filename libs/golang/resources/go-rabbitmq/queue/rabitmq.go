package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Dsn          string
	Channel      *amqp.Channel
	ExchangeName string
	ExchangeType string
}

func NewRabbitMQ(user, password, host, port, protocol, exchangeName, exchangeType string) *RabbitMQ {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/", protocol, user, password, host, port)
	return &RabbitMQ{
		Dsn:          dsn,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
	}
}

func (r *RabbitMQ) Connect() error {
	conn, err := amqp.Dial(r.Dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	r.Channel, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	return nil
}

func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		return r.Channel.Close()
	}
	return nil
}

func (r *RabbitMQ) SetupExchange() error {
	err := r.declareExchange()
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}
	return nil
}

func (r *RabbitMQ) declareExchange() error {
	return r.Channel.ExchangeDeclare(
		r.ExchangeName,
		r.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
}
