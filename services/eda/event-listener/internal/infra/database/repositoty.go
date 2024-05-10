package database

import (
	"go-knowledge/services/eda/event-listener/internal/entity"
)

type EventMessageRepositoryInterface interface {
	CreateEventMessage(eventMessage *entity.EventMessage) error
    FindAllEventMessage() ([]*entity.EventMessage, error)
    FindEventMessageByID(id string) (*entity.EventMessage, error)
}
