package database

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
)

type InMemoryDocDBTestSuite struct {
	suite.Suite
	db              *InMemoryDocBD
	dbName          string
    collectionName1 string
	collectionName2 string

}

func TestInMemoryDocDBTestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryDocDBTestSuite))
}

func (suite *InMemoryDocDBTestSuite) SetupTest() {
    suite.dbName = "test-db"
    suite.db = NewInMemoryDocBD(suite.dbName)
    suite.collectionName1 = "users"
    suite.collectionName2 = "products"
}

func (suite *InMemoryDocDBTestSuite) TearDownTest() {
    suite.db = nil
}

func (suite *InMemoryDocDBTestSuite) TestDBCreateCollection() {
    err := suite.db.CreateCollection(suite.collectionName1)
    assert.Nil(suite.T(), err)

    err = suite.db.CreateCollection(suite.collectionName2)
    assert.Nil(suite.T(), err)

    assert.Equal(suite.T(), 2, len(suite.db.Collections))
}

func (suite *InMemoryDocDBTestSuite) TestDBGetCollection() {
    err := suite.db.CreateCollection(suite.collectionName1)
    assert.Nil(suite.T(), err)

    err = suite.db.CreateCollection(suite.collectionName2)
    assert.Nil(suite.T(), err)

    collection, err := suite.db.GetCollection(suite.collectionName1)
    assert.Nil(suite.T(), err)
    assert.NotNil(suite.T(), collection)

    collection, err = suite.db.GetCollection(suite.collectionName2)
    assert.Nil(suite.T(), err)
    assert.NotNil(suite.T(), collection)
}

func (suite *InMemoryDocDBTestSuite) TestDBDropCollection() {
    err := suite.db.CreateCollection(suite.collectionName1)
    assert.Nil(suite.T(), err)

    err = suite.db.CreateCollection(suite.collectionName2)
    assert.Nil(suite.T(), err)

    err = suite.db.DropCollection(suite.collectionName1)
    assert.Nil(suite.T(), err)

    err = suite.db.DropCollection(suite.collectionName2)
    assert.Nil(suite.T(), err)

    assert.Equal(suite.T(), 0, len(suite.db.Collections))
}

func (suite *InMemoryDocDBTestSuite) TestDBListCollections() {
    err := suite.db.CreateCollection(suite.collectionName1)
    assert.Nil(suite.T(), err)

    err = suite.db.CreateCollection(suite.collectionName2)
    assert.Nil(suite.T(), err)

    collections := suite.db.ListCollections()
    assert.Equal(suite.T(), 2, len(collections))
    assert.Contains(suite.T(), collections, suite.collectionName1)
    assert.Contains(suite.T(), collections, suite.collectionName2)
}

func (suite *InMemoryDocDBTestSuite) TestDBCreateCollectionError() {
    err := suite.db.CreateCollection(suite.collectionName1)
    assert.Nil(suite.T(), err)

    err = suite.db.CreateCollection(suite.collectionName1)
    assert.NotNil(suite.T(), err)
}

func (suite *InMemoryDocDBTestSuite) TestDBDropCollectionError() {
    err := suite.db.DropCollection(suite.collectionName1)
    assert.NotNil(suite.T(), err)
}

func (suite *InMemoryDocDBTestSuite) TestDBGetCollectionError() {
    _, err := suite.db.GetCollection(suite.collectionName1)
    assert.NotNil(suite.T(), err)
}

