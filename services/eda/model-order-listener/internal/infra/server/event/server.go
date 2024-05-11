package eventserver

import (
	controllerListener "go-knowledge/services/eda/model-order-listener/internal/controller/listener"
	"log"
)

type ListenerServer struct {
	controller *controllerListener.ListenerController
	quitCh     chan struct{}
}

func NewListenerServer(
	controller *controllerListener.ListenerController,
) *ListenerServer {
	return &ListenerServer{
		controller: controller,
		quitCh:     make(chan struct{}),
	}
}

func (es *ListenerServer) Start() {
	for listenerTag, _ := range es.controller.GetListeners() {
		go es.controller.StartListener(listenerTag)
	}
mainloop:
	for {
		select {
		case <-es.quitCh:
			log.Println("Shutting down listener server")
			break mainloop
		}
	}
}
