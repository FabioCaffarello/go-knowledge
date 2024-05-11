package client

import (
	"errors"
	"go-knowledge/libs/golang/resources/database/in-memory/go-doc-db/database"
)

type Client struct {
	db *database.InMemoryDocBD
}

func NewClient(
	db *database.InMemoryDocBD,
) *Client {
	return &Client{
		db: db,
	}
}

func (c *Client) getCollection(
	collectionName string,
) (*database.Collection, error) {
	collection, err := c.db.GetCollection(collectionName)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (c *Client) CreateCollection(
	collectionName string,
) error {
	err := c.db.CreateCollection(collectionName)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DropCollection(
	collectionName string,
) error {
	err := c.db.DropCollection(collectionName)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ListCollections() []string {
	return c.db.ListCollections()
}

func (c *Client) ConvertToDocument(
	document map[string]interface{},
) (database.Document, error) {
	if document == nil {
		return nil, errors.New("document is nil")
	}
	if _, ok := document["_id"]; !ok {
		return nil, errors.New("_id field is required")
	}
	return database.Document(document), nil
}

func (c *Client) InsertOne(
	collectionName string,
	document map[string]interface{},
) error {
	collection, err := c.getCollection(collectionName)
	if err != nil {
		return err
	}
	doc, err := c.ConvertToDocument(document)
	if err != nil {
		return err
	}
	return collection.InsertOne(doc)
}

func (c *Client) FindOne(
	collectionName string,
	id string,
) (map[string]interface{}, error) {
	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}
	doc, err := collection.FindOne(id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}(doc), nil
}

func (c *Client) FindAll(
	collectionName string,
) ([]map[string]interface{}, error) {
	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}
	docs := collection.FindAll()
	documents := make([]map[string]interface{}, 0, len(docs))
	for _, doc := range docs {
		documents = append(documents, map[string]interface{}(doc))
	}
	return documents, nil
}

func (c *Client) Find(
	collectionName string,
	filter map[string]interface{},
) ([]map[string]interface{}, error) {
	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}
	docs := collection.Find(filter)
	documents := make([]map[string]interface{}, 0, len(docs))
	for _, doc := range docs {
		documents = append(documents, map[string]interface{}(doc))
	}
	return documents, nil
}

func (c *Client) UpdateOne(
	collectionName string,
	id string,
	update map[string]interface{},
) error {
	collection, err := c.getCollection(collectionName)
	if err != nil {
		return err
	}
	if update == nil || len(update) == 0 {
		return errors.New("update is empty")
	}
	return collection.UpdateOne(id, update)
}

func (c *Client) DeleteOne(
	collectionName string,
	id string,
) error {
	collection, err := c.getCollection(collectionName)
	if err != nil {
		return err
	}
	return collection.DeleteOne(id)
}

func (c *Client) DeleteAll(
	collectionName string,
) error {
	collection, err := c.getCollection(collectionName)
	if err != nil {
		return err
	}
	return collection.DeleteAll()
}
