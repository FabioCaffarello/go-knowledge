package usecase

import (
    "go-knowledge/services/eda/event-listener/internal/application/dto/input"
)

type OutputsListenerUseCase struct {

}


func NewOutputsListenerUseCase() *OutputsListenerUseCase {
    return &OutputsListenerUseCase{}
}

func (u *OutputsListenerUseCase) Execute(msg input.MessageDTO) (string, error) {
    return "", nil
}