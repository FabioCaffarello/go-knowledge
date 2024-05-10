package consumer

type ConsumerInterface interface {
    Consume()
    GetListenerTag() string
    GetMsgCh() <-chan []byte
}