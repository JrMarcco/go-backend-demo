package gorm

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type gormTestSuite struct {
	suite.Suite
}

func TestMysql(t *testing.T) {
	suite.Run(t, &gormTestSuite{})
}

func (m *gormTestSuite) SetupSuite() {
	Setup()
}
