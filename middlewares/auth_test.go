package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const authPath = "/auth"

func (m *middlewareTestSuite) buildReq(authorizationTyp, username string, duration time.Duration) *http.Request {
	t := m.T()

	req, err := http.NewRequest(http.MethodGet, authPath, nil)
	require.NoError(t, err)

	tk, err := m.server.TokenMaker.Generate(username, duration)
	require.NoError(t, err)

	authorization := fmt.Sprintf("%s %s", authorizationTyp, tk)
	req.Header.Set(authorizationKey, authorization)
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
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			// register auth middleware
			m.server.Router.GET(
				authPath,
				NewAuthMiddlewareBuilder(m.server.TokenMaker).Build(),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			m.server.Router.ServeHTTP(recorder, tc.req)

			require.Equal(t, tc.wantCode, recorder.Code)
		})
	}
}
