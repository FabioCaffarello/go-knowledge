package entity

import "errors"

type EventMessage struct {
	Data map[string]interface{} `json:"data"`
}

func NewEventMessage(data map[string]interface{}) (*EventMessage, error) {
	eventMessage := &EventMessage{
		Data: data,
	}
	if err := eventMessage.isValidate(); err != nil {
		return nil, err
	}
	return eventMessage, nil
}

func (e *EventMessage) isValidate() error {
	if e.Data == nil {
		return errors.New("data is required")
	}
	return nil
}
