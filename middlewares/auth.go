package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jrmarcco/go-backend-demo/api"
	"net/http"
	"strings"
)

const (
	authorizationKey = "authorization"
	authorizationTyp = "bearer"
	payloadKey       = "payload"
)

func AuthMiddleware(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/api/v1/user/login" {
			ctx.Next()
			return
		}

		authorizationHeader := ctx.GetHeader(authorizationKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResp(err))
			return
		}

		fds := strings.Fields(authorizationHeader)
		if len(fds) < 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResp(err))
			return
		}

		typ := strings.ToLower(fds[0])
		if typ != authorizationTyp {
			err := fmt.Errorf("unsupported authorization type: %s", typ)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResp(err))
			return
		}

		accessToken := fds[1]
		payload, err := s.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResp(err))
			return
		}

		ctx.Set(payloadKey, payload)
		ctx.Next()
	}
}
