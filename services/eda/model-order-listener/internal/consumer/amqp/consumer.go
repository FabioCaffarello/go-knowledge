package amqpconsumer

import (
	"fmt"
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpConsumer struct {
	rabbitMQ         *queue.RabbitMQ
	rabbitMQConsumer *queue.RabbitMQConsumer
	msgCh            chan []byte
	quitCh           chan struct{}
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
		quitCh:           make(chan struct{}),
	}
}

func (al *AmqpConsumer) GetListenerTag() string {
	return fmt.Sprintf("%s:%s:%s", al.rabbitMQConsumer.ConsumerName, al.rabbitMQConsumer.QueueName, al.rabbitMQConsumer.RoutingKey)
}

func (al *AmqpConsumer) Consume() {
	messageChannel := make(chan amqp.Delivery)
	go al.rabbitMQConsumer.Consume(messageChannel)
mainloop:
	for {
		select {
		case msg := <-messageChannel:
				log.Printf("Received message: %s", string(msg.Body))
				msg.Ack(false)
				al.msgCh <- msg.Body
		case <-al.quitCh:
			break mainloop
		}
	}
}

func (al *AmqpConsumer) GetMsgCh() <-chan []byte {
	return al.msgCh
}
