package main

import (
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"

	amqpServer "go-knowledge/services/eda/event-listener/internal/server"
	"log"

	"go-knowledge/services/eda/event-listener/internal/application/usecase"
)

var (
	consumerName  = "event-listener"
	exchangeName  = "services-events"
	exchangeType  = "direct"
	queueName     = "service-feedback"
	routingKey    = "services.events"
	routingKeyEOS = "services.events.eos"
)

func main() {
	rmq := queue.NewRabbitMQ("guest", "guest", "localhost", "5672", "amqp", exchangeName, exchangeType)
	if err := rmq.Connect(); err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	if err := rmq.SetupExchange(); err != nil {
		log.Fatalf("Error setting up exchange: %v", err)
	}

	notifier := queue.NewRabbitMQNotifier(rmq.Channel, rmq.ExchangeName)
	server := amqpServer.NewServer(rmq, notifier)
	ouputsListerUsecase := usecase.NewOutputsListenerUseCase()
	outputsEOSListenerUsecase := usecase.NewOutputsEOSListenerUseCase()
	server.RegisterConsumer(consumerName, exchangeName, queueName, routingKey, ouputsListerUsecase)
	server.RegisterConsumer(consumerName, exchangeName, queueName, routingKeyEOS, outputsEOSListenerUsecase)
	server.Start()
}
