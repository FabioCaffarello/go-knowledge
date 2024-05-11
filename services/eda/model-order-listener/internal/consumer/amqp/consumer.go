package amqpconsumer

import (
	"fmt"
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpConsumer struct {
	rabbitMQ         *queue.RabbitMQ
	rabbitMQConsumer *queue.RabbitMQConsumer
	msgCh            chan []byte
}

func NewAmqpConsumer(rmq *queue.RabbitMQ, queueName, consumerName, routingKey string) *AmqpConsumer {
	consumer := queue.NewRabbitMQConsumer(
		rmq.Channel,
		queueName,
		consumerName,
		rmq.ExchangeName,
		routingKey,
		false,
		nil,
	)
	return &AmqpConsumer{
		rabbitMQ:         rmq,
		rabbitMQConsumer: consumer,
		msgCh:            make(chan []byte),
	}
}

func (al *AmqpConsumer) GetListenerTag() string {
	return fmt.Sprintf("%s:%s:%s", al.rabbitMQConsumer.ConsumerName, al.rabbitMQConsumer.QueueName, al.rabbitMQConsumer.RoutingKey)
}

func (al *AmqpConsumer) Consume() {
	messageChannel := make(chan amqp.Delivery)
	go al.rabbitMQConsumer.Consume(messageChannel)
	for msg := range messageChannel {
		al.msgCh <- msg.Body
        msg.Ack(false)
	}
}

func (al *AmqpConsumer) GetMsgCh() <-chan []byte {
	return al.msgCh
}
