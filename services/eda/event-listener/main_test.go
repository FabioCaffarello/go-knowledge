package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EventListenerSuite struct {
	suite.Suite
}

func TestEventListenerSuite(t *testing.T) {
	suite.Run(t, new(EventListenerSuite))
}

func (suite *EventListenerSuite) TestHello() {
	result := FetchResource()
	assert.Equal(suite.T(), "Resources fetched", result)
}
