package usecase

import "go-knowledge/services/eda/model-order-listener/internal/types"

type UseCaseInterface interface {
    ProcessMessageChannel(msgCh <-chan []byte, errCh chan types.ErrMsg)
}
