package event

import "time"

type ErrorCreated struct {
	Name     string
	Payloads []interface{}
	Tag      string
}

func NewErrorCreated() *ErrorCreated {
	return &ErrorCreated{
		Name: "ErrorCreated",
	}
}

func (e *ErrorCreated) GetTag() string {
    return e.Tag
}

func (e *ErrorCreated) SetTag(tag string) {
    e.Tag = tag
}

func (e *ErrorCreated) GetName() string {
	return e.Name
}

func (e *ErrorCreated) GetPayloads() []interface{} {
	return e.Payloads
}

func (e *ErrorCreated) AddPayload(payload interface{}) {
	e.Payloads = append(e.Payloads, payload)
}

func (e *ErrorCreated) GetDateTime() time.Time {
	return time.Now()
}
