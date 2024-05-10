package main

import (
	inMemoryDBClient "go-knowledge/libs/golang/resources/database/in-memory/go-doc-db-client/client"
	inMemoryDB "go-knowledge/libs/golang/resources/database/in-memory/go-doc-db/database"
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
	"go-knowledge/libs/golang/services/shared/go-events/events"
	amqpConsumer "go-knowledge/services/eda/model-order-listener/internal/consumer/amqp"
	controllerListener "go-knowledge/services/eda/model-order-listener/internal/controller/listener"
	inMemoryDBRepository "go-knowledge/services/eda/model-order-listener/internal/infra/database/in-memory"
	"go-knowledge/services/eda/model-order-listener/internal/infra/event"
	eventHandler "go-knowledge/services/eda/model-order-listener/internal/infra/event/handler"
	eventServer "go-knowledge/services/eda/model-order-listener/internal/infra/server/event"
	"go-knowledge/services/eda/model-order-listener/internal/usecase"
	"log"
)

var (
	consumerName = "event-listener"
	exchangeName = "services-events"
	exchangeType = "direct"
	queueName    = "service-feedback"
	routingKey   = "services.events"
	// routingKeyEOS            = "services.events.eos"
	dbName                   = "event-messages"
	modelOrderCollectionName = "model-orders"
	rabbitmqUser             = "guest"
	rabbitmqPassword         = "guest"
	rabbitmqHost             = "localhost"
	rabbitmqPort             = "5672"
	rabbitmqProtocol         = "amqp"
)

func main() {
	db := inMemoryDB.NewInMemoryDocBD(dbName)
	dbClient := inMemoryDBClient.NewClient(db)
	modelOrderRepository := inMemoryDBRepository.NewModelOrderRepository(dbName, dbClient, modelOrderCollectionName)

	rmq := queue.NewRabbitMQ(
		rabbitmqUser,
		rabbitmqPassword,
		rabbitmqHost,
		rabbitmqPort,
		rabbitmqProtocol,
		exchangeName,
		exchangeType,
	)
	if err := rmq.Connect(); err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	if err := rmq.SetupExchange(); err != nil {
		log.Fatalf("Error setting up exchange: %v", err)
	}

	notifier := queue.NewRabbitMQNotifier(rmq.Channel, rmq.ExchangeName)
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("ErrorCreated", &eventHandler.ErrorCreatedHandler{
		Notifier: notifier,
	})

	errorEventHandler := event.NewErrorCreated()

	modelOrderUsecase := usecase.NewCreateModelOrderUseCase(
		modelOrderRepository,
		errorEventHandler,
		eventDispatcher,
	)

	listenerController := controllerListener.NewListenerController()
	modelOrderConsumer := amqpConsumer.NewAmqpConsumer(rmq, queueName, consumerName, routingKey)

	listenerController.AddListener(modelOrderConsumer, modelOrderUsecase)

	listenerServer := eventServer.NewListenerServer(listenerController)
	listenerServer.Start()
}
