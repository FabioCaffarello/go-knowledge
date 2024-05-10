package handler

import (
	"encoding/json"
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
	"go-knowledge/libs/golang/services/shared/go-events/events"
	"log"
	"sync"
)

type ErrorCreatedHandler struct {
	Notifier *queue.RabbitMQNotifier
}

func NewErrorCreatedHandler(
	notifier *queue.RabbitMQNotifier,
) *ErrorCreatedHandler {
	return &ErrorCreatedHandler{
		Notifier: notifier,
	}
}

func (h *ErrorCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range event.GetPayloads() {
		jsonOutput, _ := json.Marshal(v)
        err := h.Notifier.Notify(jsonOutput, event.GetTag())
        if err != nil {
            log.Printf("Error notifying: %v", err)
        }
	}
}
