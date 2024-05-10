package usecase

type UseCaseInterface interface {
	ProcessMessageChannel(msgCh <-chan []byte, listenerTag string)
}
