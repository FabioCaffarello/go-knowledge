package listenercontroller

import (
	"go-knowledge/services/eda/model-order-listener/internal/consumer"
	"go-knowledge/services/eda/model-order-listener/internal/usecase"
	"sync"

	"errors"
)

type Listener struct {
	consumer    consumer.ConsumerInterface
	UsecaseImpl usecase.UseCaseInterface
}

type ListenerController struct {
	listeners map[string]*Listener
	mu        sync.RWMutex
}

func NewListenerController() *ListenerController {
	return &ListenerController{}
}

func (c *ListenerController) AddListener(consum consumer.ConsumerInterface, usecaseImpl usecase.UseCaseInterface) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.listeners[consum.GetListenerTag()]
	if ok {
		return errors.New("Listener already exists")
	}
	c.listeners[consum.GetListenerTag()] = &Listener{
		consumer:    consum,
		UsecaseImpl: usecaseImpl,
	}
	return nil
}

func (c *ListenerController) RemoveListener(listenerTag string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.listeners[listenerTag]
	if !ok {
		return errors.New("Listener not found")
	}
	delete(c.listeners, listenerTag)
	return nil
}

func (c *ListenerController) GetListener(listenerTag string) (*Listener, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	listener, ok := c.listeners[listenerTag]
	if !ok {
		return &Listener{}, errors.New("Listener not found")
	}
	return listener, nil
}

func (c *ListenerController) GetListeners() map[string]*Listener {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.listeners
}

func (c *ListenerController) StartListener(listenerTag string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	listener, ok := c.listeners[listenerTag]
	if !ok {
		return errors.New("Listener not found")
	}
	go func(listener *Listener) {
		listener.consumer.Consume()
		listener.UsecaseImpl.ProcessMessageChannel(listener.consumer.GetMsgCh(), listenerTag)
	}(listener)
	return nil
}
