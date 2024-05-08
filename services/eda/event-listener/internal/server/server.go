package server

import (
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	quitCh         chan struct{}
	msgCh          chan amqp.Delivery
	rabbitMQ       *queue.RabbitMQ
	eventsNotifier *queue.RabbitMQNotifier
	consumers      []Consumer
	eventsConsumer *queue.RabbitMQConsumer
}

func NewServer(rmq *queue.RabbitMQ, notifier *queue.RabbitMQNotifier) *Server {
	return &Server{
		rabbitMQ:       rmq,
		eventsNotifier: notifier,
		quitCh:         make(chan struct{}),
		msgCh:          make(chan amqp.Delivery, 128),
	}
}

func (s *Server) RegisterConsumer(consumerName, exchangeName, queueName, routingKey string) {
	rmqConsumer := queue.NewRabbitMQConsumer(s.rabbitMQ.Channel, queueName, consumerName, exchangeName, routingKey, false, nil)
	consumer := NewConsumer(queueName, routingKey, rmqConsumer)
	s.consumers = append(s.consumers, *consumer)
}

func (s *Server) Start() {
	for _, consumer := range s.consumers {
		go consumer.consumeRabbitMQ()
	}
	s.loop()
}
