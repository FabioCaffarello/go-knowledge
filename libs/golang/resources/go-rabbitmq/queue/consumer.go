package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	Channel      *amqp.Channel
	QueueName    string
	ConsumerName string
	AutoAck      bool
	Args         amqp.Table
	ExchangeName string
	RoutingKey   string
}

func NewRabbitMQConsumer(channel *amqp.Channel, queueName, consumerName, exchangeName, routingKey string, autoAck bool, args amqp.Table) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		Channel:      channel,
		QueueName:    queueName,
		ConsumerName: consumerName,
		AutoAck:      autoAck,
		Args:         args,
		ExchangeName: exchangeName,
		RoutingKey:   routingKey,
	}
}

func (r *RabbitMQConsumer) Consume(messageChannel chan amqp.Delivery) {
	q, err := r.Channel.QueueDeclare(
		r.QueueName, // name
		true,        // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		r.Args,      // arguments
	)
	failOnError(err, "failed to declare a queue")
	err = r.Channel.QueueBind(
		q.Name,         // queue name
		r.RoutingKey,   // routing key
		r.ExchangeName, // exchange name
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "failed to bind queue to exchange")
	incomingMessage, err := r.Channel.Consume(
		q.Name,         // queue
		r.ConsumerName, // consumer
		r.AutoAck,      // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failOnError(err, "Failed to register a consumer")
	go func() {
		for message := range incomingMessage {
			log.Println("Incoming new message")
			messageChannel <- message
		}
		log.Println("RabbitMQ channel closed")
		close(messageChannel)
	}()
}

