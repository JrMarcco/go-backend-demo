package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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

func (a *apiTestSuite) newTestServer(store db.Store) *Server {
	server := NewServer(util.ServerCfg{TokenDuration: time.Minute}, store)
	server.RegisterRouter()

	return server
}
