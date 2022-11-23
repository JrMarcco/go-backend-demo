package models

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type modelTestSuite struct {
	suite.Suite
}

func TestMysql(t *testing.T) {
	suite.Run(t, &modelTestSuite{})
}

func (m *modelTestSuite) SetupSuite() {
	Setup()
}
