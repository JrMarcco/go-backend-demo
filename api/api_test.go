package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type apiTestSuite struct {
	suite.Suite
	s *S
}

func TestApi(t *testing.T) {
	suite.Run(t, &apiTestSuite{})
}

func (a *apiTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (a *apiTestSuite) setupTestServer(store db.Store) {
	server := NewServer(util.ServerCfg{TokenDuration: time.Minute}, store)
	server.RegisterRouter()

	a.s = server
}

func (a *apiTestSuite) setAuthorization(req *http.Request, username string, duration time.Duration) {
	t := a.T()

	if username != "" {
		tk, err := a.s.tokenMaker.Generate(username, duration)
		require.NoError(t, err)

		authorization := fmt.Sprintf("%s %s", authorizationTyp, tk)
		req.Header.Set(authorizationKey, authorization)
	}
}
