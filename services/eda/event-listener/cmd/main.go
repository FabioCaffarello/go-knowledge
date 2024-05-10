package main

import (
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"

	amqpServer "go-knowledge/services/eda/event-listener/internal/server"
	"log"

	"go-knowledge/services/eda/event-listener/internal/application/usecase"

	inMemoryDBClient "go-knowledge/libs/golang/resources/database/in-memory/go-doc-db-client/client"
	inMemoryDB "go-knowledge/libs/golang/resources/database/in-memory/go-doc-db/database"
	inMemoryDBRepository "go-knowledge/services/eda/event-listener/internal/infra/database/in-memory"
)

var (
	consumerName  = "event-listener"
	exchangeName  = "services-events"
	exchangeType  = "direct"
	queueName     = "service-feedback"
	routingKey    = "services.events"
	routingKeyEOS = "services.events.eos"
    dbName = "event-messages"
)

func main() {
	db := inMemoryDB.NewInMemoryDocBD(dbName)
	dbClient := inMemoryDBClient.NewClient(db)
	eventsMessageRepository := inMemoryDBRepository.NewEventMessageRepository(
		log.Default(),
		"event-listener",
		dbClient,
        dbName,
	)

	rmq := queue.NewRabbitMQ("guest", "guest", "localhost", "5672", "amqp", exchangeName, exchangeType)
	if err := rmq.Connect(); err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	if err := rmq.SetupExchange(); err != nil {
		log.Fatalf("Error setting up exchange: %v", err)
	}

	notifier := queue.NewRabbitMQNotifier(rmq.Channel, rmq.ExchangeName)
	server := amqpServer.NewServer(rmq, notifier)

	ouputsListerUsecase := usecase.NewOutputsListenerUseCase(eventsMessageRepository, notifier)
	outputsEOSListenerUsecase := usecase.NewOutputsEOSListenerUseCase(eventsMessageRepository, notifier)

	server.RegisterConsumer(consumerName, exchangeName, queueName, routingKey, ouputsListerUsecase)
	server.RegisterConsumer(consumerName, exchangeName, queueName, routingKeyEOS, outputsEOSListenerUsecase)
	server.Start()
}
