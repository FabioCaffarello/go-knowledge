package event

import "time"

type FileOrderCreated struct {
	Name     string
	Payloads []interface{}
	Tag      string
}

func NewFileOrderCreated() *FileOrderCreated {
	return &FileOrderCreated{
		Name: "FileOrderCreated",
	}
}

func (e *FileOrderCreated) GetTag() string {
    return e.Tag
}

func (e *FileOrderCreated) SetTag(tag string) {
    e.Tag = tag
}

func (e *FileOrderCreated) GetName() string {
	return e.Name
}

func (e *FileOrderCreated) GetPayloads() []interface{} {
	return e.Payloads
}

func (e *FileOrderCreated) AddPayload(payload interface{}) {
	e.Payloads = append(e.Payloads, payload)
}

func (e *FileOrderCreated) GetDateTime() time.Time {
	return time.Now()
}
