package server

import (
	"fmt"
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ConsumerTag      string
	QueueName        string
	RoutingKey       string
	rabbitMQConsumer *queue.RabbitMQConsumer
	quitCh           chan struct{}
	msgCh            chan amqp.Delivery
}

func NewConsumer(queueName, routingKey string, rmqConsumer *queue.RabbitMQConsumer) *Consumer {
	return &Consumer{
		QueueName:        queueName,
		RoutingKey:       routingKey,
		rabbitMQConsumer: rmqConsumer,
		quitCh:           make(chan struct{}),
		msgCh:            make(chan amqp.Delivery, 128),
	}
}

func (c *Consumer) GetConsumerTag() string {
	return fmt.Sprintf("%s:%s:%s", c.rabbitMQConsumer.ConsumerName, c.rabbitMQConsumer.QueueName, c.RoutingKey)
}

func (c *Consumer) consumeRabbitMQ() {
	go c.rabbitMQConsumer.Consume(c.msgCh)
	for {
		select {
		case msg, ok := <-c.msgCh:
			if !ok {
				log.Println("Message channel closed, stopping consumption")
				return
			}
			log.Printf("Received message: %s", string(msg.Body))
			msg.Ack(false)
		default:
		}
	}
}
