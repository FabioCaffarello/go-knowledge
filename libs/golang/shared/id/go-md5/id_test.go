package gomd5id

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GoMd5IDTestSuite struct {
	suite.Suite
}

func TestGoMd5IDTestSuite(t *testing.T) {
	suite.Run(t, new(GoMd5IDTestSuite))
}

func (suite *GoMd5IDTestSuite) TestGetIDFromString() {
	id := GetIDFromString("hello")
	assert.Equal(suite.T(), "5d41402abc4b2a76b9719d911017c592", string(id))
}
