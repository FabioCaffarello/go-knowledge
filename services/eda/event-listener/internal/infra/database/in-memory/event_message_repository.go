package inmemorydb

import (
	"encoding/json"
	"log"

	"go-knowledge/libs/golang/resources/database/in-memory/go-doc-db-client/client"
	"go-knowledge/services/eda/event-listener/internal/entity"
)

type EventMessageRepository struct {
	log            *log.Logger
	databaseName   string
	client         *client.Client
	collectionName string
}

func NewEventMessageRepository(
	log *log.Logger,
	databaseName string,
	client *client.Client,
	collectionName string,
) *EventMessageRepository {
	repository := &EventMessageRepository{
		log:            log,
		databaseName:   databaseName,
		client:         client,
		collectionName: collectionName,
	}
	repository.client.CreateCollection(collectionName)
	return repository
}

func (r *EventMessageRepository) CreateEventMessage(eventMessage *entity.EventMessage) error {
    r.log.Println("Save event message in memory")
    eventMessageMap := map[string]interface{}{
        "data": eventMessage.Data,
    }
    err := r.client.InsertOne(r.collectionName, eventMessageMap)
    if err != nil {
        return err
    }
    return nil
}

func (r *EventMessageRepository) FindAllEventMessage() ([]*entity.EventMessage, error) {
	r.log.Println("Find all event messages in memory")
	documents, err := r.client.FindAll(r.collectionName)
	if err != nil {
		return nil, err
	}
	eventMessages := make([]*entity.EventMessage, 0, len(documents))
	for _, document := range documents {

		var result entity.EventMessage
		documentBytes, err := json.Marshal(document)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(documentBytes, &result)
		if err != nil {
			return nil, err
		}

		eventMessages = append(eventMessages, &result)
	}
	return eventMessages, nil
}

func (r *EventMessageRepository) FindEventMessageByID(id string) (*entity.EventMessage, error) {
	r.log.Println("Find event message by ID in memory")
	document, err := r.client.FindOne(r.collectionName, id)
	if err != nil {
		return nil, err
	}
	var result entity.EventMessage
	documentBytes, err := json.Marshal(document)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(documentBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
