package eventserver

import (
	controllerListener "go-knowledge/services/eda/model-order-listener/internal/controller/listener"
)

type ListenerServer struct {
	controller *controllerListener.ListenerController
}

func NewListenerServer(
	controller *controllerListener.ListenerController,
) *ListenerServer {
	return &ListenerServer{
		controller: controller,
	}
}

func (es *ListenerServer) Start() {
	for listenerTag, _ := range es.controller.GetListeners() {
		es.controller.StartListener(listenerTag)
	}
mainloop:
	for {
		select {
		case <-es.controller.GetQuitCh():
			break mainloop
		}
	}
}
