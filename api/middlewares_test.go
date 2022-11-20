package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type middlewareTestSuite struct {
	suite.Suite
	s *S
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, &middlewareTestSuite{
		s: NewServer(util.ServerCfg{}, nil),
	})
}

func (m *middlewareTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

const authPath = "/auth"

func (m *middlewareTestSuite) buildReq(authorizationTyp, username string, duration time.Duration) *http.Request {
	t := m.T()

	req, err := http.NewRequest(http.MethodGet, authPath, nil)
	require.NoError(t, err)

	if username != "" {
		tk, err := m.s.tokenMaker.Generate(username, duration)
		require.NoError(t, err)

		authorization := fmt.Sprintf("%s %s", authorizationTyp, tk)
		req.Header.Set(authorizationKey, authorization)
	}

	return req
}

func (m *middlewareTestSuite) TestAuthMiddleware() {
	t := m.T()

	tcs := []struct {
		name     string
		req      *http.Request
		wantCode int
	}{
		{
			name:     "Normal Case",
			req:      m.buildReq(authorizationTyp, "jrmarcco", time.Minute),
			wantCode: http.StatusOK,
		},
		{
			name:     "Unauthorized Case",
			req:      m.buildReq(authorizationTyp, "", time.Minute),
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Unsupported Case",
			req:      m.buildReq("unsupported", "jrmarcco", time.Minute),
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Invalid Format Case",
			req:      m.buildReq("", "jrmarcco", time.Minute),
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Expired Case",
			req:      m.buildReq(authorizationTyp, "jrmarcco", -time.Minute),
			wantCode: http.StatusUnauthorized,
		},
	}

	// register auth middleware
	m.s.router.GET(
		authPath,
		NewAuthMiddlewareBuilder(m.s.tokenMaker).Build(),
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
		},
	)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			recorder := httptest.NewRecorder()
			m.s.router.ServeHTTP(recorder, tc.req)

			require.Equal(t, tc.wantCode, recorder.Code)
		})
	}
}
