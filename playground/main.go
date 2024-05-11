package main

import (
	"encoding/json"
	"fmt"
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	quitCh         chan struct{}
	msgCh          chan amqp.Delivery
	tempCh         chan string
	rabbitMQ       *queue.RabbitMQ
	eventsNotifier *queue.RabbitMQNotifier
	eventsConsumer *queue.RabbitMQConsumer
}

func NewServer(rmq *queue.RabbitMQ, notifier *queue.RabbitMQNotifier, consumer *queue.RabbitMQConsumer) *Server {
	return &Server{
		rabbitMQ:       rmq,
		eventsNotifier: notifier,
		eventsConsumer: consumer,
		quitCh:         make(chan struct{}),
		msgCh:          make(chan amqp.Delivery, 128),
		tempCh:         make(chan string, 128),
	}
}

func (s *Server) Start() {
	fmt.Println("Server starting...")
	go s.consumeRabbitMQ()
	s.loop()
	defer s.rabbitMQ.Close()
}

func (s *Server) consumeRabbitMQ() {
	go s.eventsConsumer.Consume(s.msgCh)
	for {
		select {
		case msg, ok := <-s.msgCh:
			if !ok {
				log.Println("Message channel closed, stopping consumption")
				return
			}
			msg.Ack(false)
			s.tempCh <- string(msg.Body)
		default:
		}
	}
}

func (s *Server) cleanupRabbitMQ() {
	// Perform cleanup, e.g., ack messages, close channels, etc.
	if s.rabbitMQ != nil {
		s.rabbitMQ.Close()
	}
}

func (s *Server) loop() {
mainloop:
	for {
		select {
		case <-s.quitCh:
			fmt.Println("Quitting server...")
			s.cleanupRabbitMQ()
			break mainloop
		case msg := <-s.tempCh:
			s.handleMessage(msg)
		}
	}
}

func (s *Server) handleMessage(msg string) {
	fmt.Printf("Handling message: %v\n", msg)
}

func (s *Server) quit() {
	s.quitCh <- struct{}{}
}

type SubcontextDTO struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

type FileDTO struct {
	Name string `json:"name"`
}

type ModelOrderDTO struct {
	Costumer    string          `json:"costumer"`
	Context     string          `json:"context"`
	Subcontexts []SubcontextDTO `json:"subcontexts"`
	BucketName  string          `json:"bucket_name"`
	Files       []FileDTO       `json:"files"`
	Partition   string          `json:"partition"`
}

var (
	consumerName             = "event-listener"
	exchangeName             = "services-events"
	exchangeType             = "direct"
	queueName                = "service-feedback"
	routingKey               = "services.events"
	modelOrderCollectionName = "model-orders"
	rabbitmqUser             = "guest"
	rabbitmqPassword         = "guest"
	rabbitmqHost             = "localhost"
	rabbitmqPort             = "5672"
	rabbitmqProtocol         = "amqp"
	modelOrderDTO            = ModelOrderDTO{
		Costumer: "costumer",
		Context:  "context",
		Subcontexts: []SubcontextDTO{
			{
				Name:     "subcontext",
				Priority: 1,
			},
			{
				Name:     "subcontext2",
				Priority: 2,
			},
		},
		BucketName: "bucket_name",
		Files: []FileDTO{
			{
				Name: "file",
			},
			{
				Name: "file2",
			},
		},
		Partition: "partition",
	}
)

func main() {
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
	// consumer := queue.NewRabbitMQConsumer(rmq.Channel, "service-feedback", "playground-consumer", "services-events", "services.events", false, nil)

	// server := NewServer(rmq, notifier, consumer)
	// // go func() {
	// // 	time.Sleep(20 * time.Second)
	// // 	server.quit()
	// // }()
	// server.Start()

	data, err := json.Marshal(modelOrderDTO)
	if err != nil {
		log.Fatalf("Error marshaling modelOrderDTO: %v", err)
	}
	notifier.Notify(data, routingKey)
}
