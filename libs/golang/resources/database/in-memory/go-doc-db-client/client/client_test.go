package client

import (
	"go-knowledge/libs/golang/resources/database/in-memory/go-doc-db/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InMemoryDocDBClientTestSuite struct {
	suite.Suite
	db              *database.InMemoryDocBD
	client          *Client
	collectionName1 string
	collectionName2 string
	document1       map[string]interface{}
	document2       map[string]interface{}
}

func TestInMemoryDocDBClientTestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryDocDBClientTestSuite))
}

func (suite *InMemoryDocDBClientTestSuite) SetupTest() {
	suite.db = database.NewInMemoryDocBD("test-db")
	suite.client = NewClient(suite.db)
	suite.collectionName1 = "users"
	suite.collectionName2 = "products"
	suite.document1 = map[string]interface{}{
		"_id":  "1",
		"name": "Alice",
		"age":  30,
	}
	suite.document2 = map[string]interface{}{
		"_id":  "2",
		"name": "Bob",
		"age":  25,
	}
}

func (suite *InMemoryDocDBClientTestSuite) TearDownTest() {
	suite.db = nil
	suite.client = nil
	suite.document1 = nil
	suite.document2 = nil
}

func (suite *InMemoryDocDBClientTestSuite) TestClientCreateCollection() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.CreateCollection(suite.collectionName2)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), 2, len(suite.db.Collections))
}

func (suite *InMemoryDocDBClientTestSuite) TestClientGetCollection() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.CreateCollection(suite.collectionName2)
	assert.Nil(suite.T(), err)

	collection, err := suite.client.getCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), collection)

	collection, err = suite.client.getCollection(suite.collectionName2)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), collection)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientDropCollection() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.CreateCollection(suite.collectionName2)
	assert.Nil(suite.T(), err)

	err = suite.client.DropCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.DropCollection(suite.collectionName2)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), 0, len(suite.db.Collections))
}

func (suite *InMemoryDocDBClientTestSuite) TestClientListCollections() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.CreateCollection(suite.collectionName2)
	assert.Nil(suite.T(), err)

	collections := suite.client.ListCollections()
	assert.Equal(suite.T(), 2, len(collections))
	assert.Contains(suite.T(), collections, suite.collectionName1)
	assert.Contains(suite.T(), collections, suite.collectionName2)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientInsertOne() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document2)
	assert.Nil(suite.T(), err)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientInsertOneError() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, nil)
	assert.NotNil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.NotNil(suite.T(), err)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientFindOne() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	document, err := suite.client.FindOne(suite.collectionName1, "1")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.document1, document)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientFindOneError() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	document, err := suite.client.FindOne(suite.collectionName1, "2")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), document)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientFindAll() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document2)
	assert.Nil(suite.T(), err)

	documents, err := suite.client.FindAll(suite.collectionName1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, len(documents))
	assert.Contains(suite.T(), documents, suite.document1)
	assert.Contains(suite.T(), documents, suite.document2)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientFindAllError() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	documents, err := suite.client.FindAll(suite.collectionName1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(documents))
	assert.Contains(suite.T(), documents, suite.document1)

	documents, err = suite.client.FindAll(suite.collectionName2)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), documents)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientFind() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document2)
	assert.Nil(suite.T(), err)

	filter := map[string]interface{}{
		"age": 30,
	}
	documents, err := suite.client.Find(suite.collectionName1, filter)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(documents))
	assert.Contains(suite.T(), documents, suite.document1)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientFindError() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	filter := map[string]interface{}{
		"age": 30,
	}
	documents, err := suite.client.Find(suite.collectionName1, filter)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(documents))
	assert.Contains(suite.T(), documents, suite.document1)

	filter = map[string]interface{}{
		"age": 25,
	}
	documents, err = suite.client.Find(suite.collectionName1, filter)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, len(documents))
}

func (suite *InMemoryDocDBClientTestSuite) TestClientFindWithNestedFields() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	document := map[string]interface{}{
		"_id": "1",
		"person": map[string]interface{}{
			"name": "Charlie",
			"age":  35,
		},
	}

	err = suite.client.InsertOne(suite.collectionName1, document)
	assert.Nil(suite.T(), err)

	query := map[string]interface{}{
		"person": map[string]interface{}{
			"name": "Charlie",
		},
	}
	documents, err := suite.client.Find(suite.collectionName1, query)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(documents))
	assert.Contains(suite.T(), documents, document)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientConvertToDocument() {
	document := map[string]interface{}{
		"_id":  "1",
		"name": "Alice",
		"age":  30,
	}
	doc, err := suite.client.ConvertToDocument(document)
	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), database.Document{}, doc)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientConvertToDocumentError() {
	doc, err := suite.client.ConvertToDocument(nil)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), doc)

	doc2 := map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}
	doc, err = suite.client.ConvertToDocument(doc2)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), doc)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientDeleteOne() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.client.DeleteOne(suite.collectionName1, "1")
	assert.Nil(suite.T(), err)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientDeleteOneError() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.client.DeleteOne(suite.collectionName1, "2")
	assert.NotNil(suite.T(), err)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientUpdateOne() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	update := map[string]interface{}{
		"age": 31,
	}
	err = suite.client.UpdateOne(suite.collectionName1, "1", update)
	assert.Nil(suite.T(), err)

	document, err := suite.client.FindOne(suite.collectionName1, "1")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 31, document["age"])
}

func (suite *InMemoryDocDBClientTestSuite) TestClientUpdateOneError() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	update := map[string]interface{}{
		"age": 31,
	}
	err = suite.client.UpdateOne(suite.collectionName1, "2", update)
	assert.NotNil(suite.T(), err)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientUpdateOneInvalidUpdate() {
	err := suite.client.CreateCollection(suite.collectionName1)
	assert.Nil(suite.T(), err)

	err = suite.client.InsertOne(suite.collectionName1, suite.document1)
	assert.Nil(suite.T(), err)

	update := map[string]interface{}{
		"age": 31,
	}
	err = suite.client.UpdateOne(suite.collectionName2, "1", update)
	assert.NotNil(suite.T(), err)

    err = suite.client.UpdateOne(suite.collectionName1, "1", nil)
    assert.NotNil(suite.T(), err)

    err = suite.client.UpdateOne(suite.collectionName1, "1", map[string]interface{}{})
    assert.NotNil(suite.T(), err)
}

func (suite *InMemoryDocDBClientTestSuite) TestClientDeleteAll() {
    err := suite.client.CreateCollection(suite.collectionName1)
    assert.Nil(suite.T(), err)

    err = suite.client.InsertOne(suite.collectionName1, suite.document1)
    assert.Nil(suite.T(), err)

    err = suite.client.InsertOne(suite.collectionName1, suite.document2)
    assert.Nil(suite.T(), err)

    err = suite.client.DeleteAll(suite.collectionName1)
    assert.Nil(suite.T(), err)

    documents, err := suite.client.FindAll(suite.collectionName1)
    assert.Nil(suite.T(), err)
    assert.Equal(suite.T(), 0, len(documents))
}

func (suite *InMemoryDocDBClientTestSuite) TestClientDeleteAllError() {
    err := suite.client.CreateCollection(suite.collectionName1)
    assert.Nil(suite.T(), err)

    err = suite.client.InsertOne(suite.collectionName1, suite.document1)
    assert.Nil(suite.T(), err)

    err = suite.client.DeleteAll(suite.collectionName2)
    assert.NotNil(suite.T(), err)
}


