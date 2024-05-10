package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CollectionTestSuite struct {
	suite.Suite
	collection *Collection
	document1  Document
	document2  Document
}

func TestCollectionTestSuite(t *testing.T) {
	suite.Run(t, new(CollectionTestSuite))
}

func (suite *CollectionTestSuite) SetupTest() {
	suite.collection = NewCollection()
	suite.document1 = Document{
		"_id":  "1",
		"name": "Alice",
		"age":  30,
	}
	suite.document2 = Document{
		"_id":  "2",
		"name": "Bob",
		"age":  25,
	}
}

func (suite *CollectionTestSuite) TearDownTest() {
    suite.collection = nil
    suite.document1 = nil
    suite.document2 = nil
}

func (suite *CollectionTestSuite) TestCollectionInsert() {
	err := suite.collection.InsertOne(suite.document1)
	suite.Nil(err)

	err = suite.collection.InsertOne(suite.document2)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), 2, len(suite.collection.data))
}

func (suite *CollectionTestSuite) TestCollectionFindOne() {
	err := suite.collection.InsertOne(suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.collection.InsertOne(suite.document2)
	assert.Nil(suite.T(), err)

	document, err := suite.collection.FindOne("1")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.document1, document)
}

func (suite *CollectionTestSuite) TestCollectionFindAll() {
	err := suite.collection.InsertOne(suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.collection.InsertOne(suite.document2)
	assert.Nil(suite.T(), err)

	documents := suite.collection.FindAll()
	assert.Equal(suite.T(), 2, len(documents))
	assert.Contains(suite.T(), documents, suite.document1)
	assert.Contains(suite.T(), documents, suite.document2)
}

func (suite *CollectionTestSuite) TestCollectionFind() {
	err := suite.collection.InsertOne(suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.collection.InsertOne(suite.document2)
	assert.Nil(suite.T(), err)

	documents := suite.collection.Find(map[string]interface{}{"name": "Alice"})
	assert.Equal(suite.T(), 1, len(documents))
	assert.Contains(suite.T(), documents, suite.document1)
}

func (suite *CollectionTestSuite) TestCollectionFindWithNested() {
	document := Document{
		"_id": "3",
		"person": map[string]interface{}{
			"name": "Charlie",
			"age":  35,
		},
	}

	err := suite.collection.InsertOne(document)
	assert.Nil(suite.T(), err)

	query := map[string]interface{}{
		"person": map[string]interface{}{
			"name": "Charlie",
		},
	}
	documents := suite.collection.Find(query)
	assert.Equal(suite.T(), 1, len(documents))
	assert.Contains(suite.T(), documents, document)
}

func (suite *CollectionTestSuite) TestCollectionFindNone() {
	err := suite.collection.InsertOne(suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.collection.InsertOne(suite.document2)
	assert.Nil(suite.T(), err)

	documents := suite.collection.Find(map[string]interface{}{"name": "Charlie"})
	assert.Equal(suite.T(), 0, len(documents))
}

func (suite *CollectionTestSuite) TestCollectionFindAllEmpty() {
	documents := suite.collection.FindAll()
	assert.Equal(suite.T(), 0, len(documents))
}

func (suite *CollectionTestSuite) TestCollectionFindEmpty() {
	documents := suite.collection.Find(map[string]interface{}{"name": "Charlie"})
	assert.Equal(suite.T(), 0, len(documents))
}

func (suite *CollectionTestSuite) TestCollectionInsertDuplicate() {
	err := suite.collection.InsertOne(suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.collection.InsertOne(suite.document1)
	assert.NotNil(suite.T(), err)
}

func (suite *CollectionTestSuite) TestCollectionFindOneNotFound() {
	_, err := suite.collection.FindOne("1")
	assert.NotNil(suite.T(), err)
}

func (suite *CollectionTestSuite) TestCollectionFindOneEmpty() {
	document, err := suite.collection.FindOne("1")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), document)
}

func (suite *CollectionTestSuite) TestCollectionFindOneInvalidID() {
	document, err := suite.collection.FindOne("invalid")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), document)
}

func (suite *CollectionTestSuite) TestCollectionFindNoneEmpty() {
	documents := suite.collection.Find(map[string]interface{}{"name": "Alice"})
	assert.Equal(suite.T(), 0, len(documents))
}

func (suite *CollectionTestSuite) TestCollectionFindInvalidQuery() {
	documents := suite.collection.Find(map[string]interface{}{"invalid": "Alice"})
	assert.Equal(suite.T(), 0, len(documents))
}

func (suite *CollectionTestSuite) TestCollectionFindInvalidQueryType() {
	documents := suite.collection.Find(map[string]interface{}{"age": "Alice"})
	assert.Equal(suite.T(), 0, len(documents))
}

func (suite *CollectionTestSuite) TestCollectionFindInvalidQueryValue() {
	documents := suite.collection.Find(map[string]interface{}{"age": 30})
	assert.Equal(suite.T(), 0, len(documents))
}

func (suite *CollectionTestSuite) TestDeleteOne() {
	err := suite.collection.InsertOne(suite.document1)
	assert.Nil(suite.T(), err)

	// Test deleting the document
	err = suite.collection.DeleteOne(suite.document1["_id"].(string))
	assert.Nil(suite.T(), err)

	// Verify the document is deleted
	_, err = suite.collection.FindOne(suite.document1["_id"].(string))
	assert.NotNil(suite.T(), err)

	// Attempt to delete non-existing document
	err = suite.collection.DeleteOne("non_existent_id")
	assert.Equal(suite.T(), errors.New("document not found"), err)
}

func (suite *CollectionTestSuite) TestDeleteOneEmpty() {
	err := suite.collection.DeleteOne("1")
	assert.NotNil(suite.T(), err)
}

func (suite *CollectionTestSuite) TestDeleteOneInvalidID() {
	err := suite.collection.DeleteOne("invalid")
	assert.NotNil(suite.T(), err)
}

func (suite *CollectionTestSuite) TestUpdateOneValid() {
	document1ID := suite.document1["_id"].(string)
	err := suite.collection.InsertOne(suite.document1)
	assert.Nil(suite.T(), err)

	// Update the document
	updatedDocument := Document{
		"_id":  "1",
		"name": "Alice",
		"age":  31,
	}
	err = suite.collection.UpdateOne(document1ID, updatedDocument)
	assert.Nil(suite.T(), err)

	// Verify the document is updated
	document, err := suite.collection.FindOne(document1ID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), updatedDocument, document)
}

func (suite *CollectionTestSuite) TestUpdateOneEmpty() {
	err := suite.collection.UpdateOne("1", suite.document1)
	assert.NotNil(suite.T(), err)
}

func (suite *CollectionTestSuite) TestDeleteAll() {
	err := suite.collection.InsertOne(suite.document1)
	assert.Nil(suite.T(), err)

	err = suite.collection.InsertOne(suite.document2)
	assert.Nil(suite.T(), err)

	// Test deleting all documents
	err = suite.collection.DeleteAll()
	assert.Nil(suite.T(), err)

	// Verify all documents are deleted
	documents := suite.collection.FindAll()
	assert.Equal(suite.T(), 0, len(documents))
}
