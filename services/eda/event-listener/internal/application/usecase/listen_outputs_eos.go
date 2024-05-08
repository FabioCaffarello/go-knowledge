package usecase

import (
	"go-knowledge/services/eda/event-listener/internal/application/dto/input"
	"go-knowledge/services/eda/event-listener/internal/entity"
	"log"
)

type OutputsEOSListenerUseCase struct {
}

func NewOutputsEOSListenerUseCase() *OutputsEOSListenerUseCase {
    return &OutputsEOSListenerUseCase{}
}

func (u *OutputsEOSListenerUseCase) Execute(msg input.MessageDTO) (string, error) {
    msgEntity, err := entity.NewEventMessage(msg.Data)
    if err != nil {
        log.Printf("Error creating event message: %v", err)
        return "", err
    }

    // store in database in memory

    log.Printf("Received message: %+v", msgEntity)
    return "", nil
}