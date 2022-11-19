package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jrmarcco/go-backend-demo/api"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/suite"
	"testing"
)

type middlewareTestSuite struct {
	suite.Suite
	server *api.Server
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, &middlewareTestSuite{
		server: api.NewServer(util.ServerCfg{}, nil),
	})
}

func (m *middlewareTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}
