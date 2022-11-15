package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"testing"
)

type apiTestSuite struct {
	suite.Suite
}

func TestApi(t *testing.T) {
	suite.Run(t, &apiTestSuite{})
}

func (a *apiTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}
