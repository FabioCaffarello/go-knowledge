package usecase

import (
    "go-knowledge/services/eda/event-listener/internal/application/dto/input"
)

type UsecaseInterface interface {
	Execute(msg input.MessageDTO) (string, error)
}