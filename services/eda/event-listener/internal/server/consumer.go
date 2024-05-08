package server

import (
	"encoding/json"
	"fmt"
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
	"log"
	"sync"

	inputDTO "go-knowledge/services/eda/event-listener/internal/application/dto/input"
	"go-knowledge/services/eda/event-listener/internal/application/usecase"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ConsumerTag      string
	QueueName        string
	RoutingKey       string
	rabbitMQConsumer *queue.RabbitMQConsumer
	quitCh           chan struct{}
	msgCh            chan amqp.Delivery
	serverQuitCh     chan string
    usecaseImpl      usecase.UsecaseInterface
}

func NewConsumer(queueName, routingKey string, rmqConsumer *queue.RabbitMQConsumer, serverQuitCh chan string, usecaseImpl usecase.UsecaseInterface) *Consumer {
	return &Consumer{
		QueueName:        queueName,
		RoutingKey:       routingKey,
		rabbitMQConsumer: rmqConsumer,
		quitCh:           make(chan struct{}),
		msgCh:            make(chan amqp.Delivery, 128),
        serverQuitCh:    serverQuitCh,
        usecaseImpl:     usecaseImpl,
	}
}

func (c *Consumer) GetConsumerTag() string {
	return fmt.Sprintf("%s:%s:%s", c.rabbitMQConsumer.ConsumerName, c.rabbitMQConsumer.QueueName, c.RoutingKey)
}

func (c *Consumer) processMessage(wg sync.WaitGroup, body []byte) error {
    var message inputDTO.MessageDTO
    err := json.Unmarshal(body, &message)
    if err != nil {
        log.Printf("Error unmarshalling message: %v", err)
        return err
    }
	log.Printf("Received message: %+v", message)
    // usecase
    c.usecaseImpl.Execute(message)

	wg.Done()
    return nil
}

func (c *Consumer) Listen() {
	go c.rabbitMQConsumer.Consume(c.msgCh)
mainloop:
	for {
		select {
		case <-c.quitCh:
			fmt.Println("Stopping consumer...")
			break mainloop
		case msg, ok := <-c.msgCh:
			if !ok {
				log.Println("Message channel closed, stopping consumption")
				return
			}
			var wg sync.WaitGroup
			wg.Add(1)
			err := c.processMessage(wg, msg.Body)
            if err != nil {
                c.serverQuitCh <- c.GetConsumerTag()
            }
			wg.Wait()
			msg.Ack(false)
		}
	}
}
