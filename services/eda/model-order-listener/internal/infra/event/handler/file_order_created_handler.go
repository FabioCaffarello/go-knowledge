package handler

import (
	"encoding/json"
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
	"go-knowledge/libs/golang/services/shared/go-events/events"
	"log"
	"sync"
)

type FileOrderCreatedHandler struct {
	Notifier *queue.RabbitMQNotifier
}

func NewFileOrderCreatedHandler(
	notifier *queue.RabbitMQNotifier,
) *FileOrderCreatedHandler {
	return &FileOrderCreatedHandler{
		Notifier: notifier,
	}
}

func (h *FileOrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range event.GetPayloads() {
		jsonOutput, err := json.Marshal(v)
		if err != nil {
			log.Printf("FileOrder marshalling: %v", err)
		}
		log.Printf("FileOrder created: %s", string(jsonOutput))
		log.Printf("FileOrder created: %s", event.GetTag())
		err = h.Notifier.Notify(jsonOutput, event.GetTag())
		if err != nil {
			log.Printf("FileOrder notifying: %v", err)
		}
	}
}
