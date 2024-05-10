package usecase

import (
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
	"go-knowledge/services/eda/event-listener/internal/application/dto/input"
	"go-knowledge/services/eda/event-listener/internal/entity"
	repository "go-knowledge/services/eda/event-listener/internal/infra/database"
	"log"
)

type OutputsEOSListenerUseCase struct {
    eventMessageRepository repository.EventMessageRepositoryInterface
    eventsNotifier         *queue.RabbitMQNotifier
}

func NewOutputsEOSListenerUseCase(
    eventMessageRepository repository.EventMessageRepositoryInterface,
    eventsNotifier *queue.RabbitMQNotifier,
) *OutputsEOSListenerUseCase {
    return &OutputsEOSListenerUseCase{
        eventMessageRepository: eventMessageRepository,
        eventsNotifier:         eventsNotifier,
    }
}

func (u *OutputsEOSListenerUseCase) Execute(msg input.MessageDTO) (string, error) {
    msgEntity, err := entity.NewEventMessage(msg.Data)
    if err != nil {
        log.Printf("Error creating event message: %v", err)
        return "", err
    }

    // store in database in memory
    err = u.eventMessageRepository.CreateEventMessage(msgEntity)
    if err != nil {
        log.Printf("Error storing event message: %v", err)
        return "", err
    }

    log.Printf("Received message: %+v", msgEntity)
    return "", nil
}