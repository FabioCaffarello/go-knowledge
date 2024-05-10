package usecase

type UseCaseInterface interface {
    ProcessMessageChannel(msgCh <-chan []byte, quitCh chan struct{})
}