package server

import (
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
    "go-knowledge/services/eda/event-listener/internal/application/usecase"
)

type Server struct {
	quitCh         chan string
	rabbitMQ       *queue.RabbitMQ
	eventsNotifier *queue.RabbitMQNotifier
	consumers      []Consumer
}

func NewServer(rmq *queue.RabbitMQ, notifier *queue.RabbitMQNotifier) *Server {
	return &Server{
		rabbitMQ:       rmq,
		eventsNotifier: notifier,
		quitCh:         make(chan string),
	}
}

func (s *Server) RegisterConsumer(consumerName, exchangeName, queueName, routingKey string, usecaseImpl usecase.UsecaseInterface) {
	rmqConsumer := queue.NewRabbitMQConsumer(s.rabbitMQ.Channel, queueName, consumerName, exchangeName, routingKey, false, nil)
	consumer := NewConsumer(queueName, routingKey, rmqConsumer, s.quitCh, usecaseImpl)
	s.consumers = append(s.consumers, *consumer)
}

func (s *Server) Start() {
	for _, consumer := range s.consumers {
		go consumer.Listen()
	}
// mainloop:
// 	for {
// 		select {
// 		case <-s.quitCh:
// 			// remove consumer
// 			// break mainloop
// 		}
// 	}
}
