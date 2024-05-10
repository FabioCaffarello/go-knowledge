package database

import "errors"

type InMemoryDocBD struct {
	Name        string
	Collections map[string]*Collection
}

func NewInMemoryDocBD(name string) *InMemoryDocBD {
	return &InMemoryDocBD{
		Name:        name,
		Collections: make(map[string]*Collection),
	}
}

func (d *InMemoryDocBD) GetCollection(
	collectionName string,
) (*Collection, error) {
	collection, ok := d.Collections[collectionName]
	if !ok {
		return nil, errors.New("collection not found")
	}
	return collection, nil
}

func (d *InMemoryDocBD) CreateCollection(
	collectionName string,
) error {
	if _, ok := d.Collections[collectionName]; ok {
		return errors.New("collection already exists")
	}
	d.Collections[collectionName] = NewCollection()
	return nil
}

func (d *InMemoryDocBD) DropCollection(
	collectionName string,
) error {
	if _, ok := d.Collections[collectionName]; !ok {
		return errors.New("collection not found")
	}
	delete(d.Collections, collectionName)
	return nil
}

func (d *InMemoryDocBD) ListCollections() []string {
	collectionNames := make([]string, 0, len(d.Collections))
	for collectionName := range d.Collections {
		collectionNames = append(collectionNames, collectionName)
	}
	return collectionNames
}
