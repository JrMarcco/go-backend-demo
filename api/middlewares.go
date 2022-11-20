package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jrmarcco/go-backend-demo/token"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const (
	authorizationKey = "authorization"
	authorizationTyp = "bearer"
	payloadKey       = "payload"
)

type AuthMiddlewareBuilder struct {
	maker token.Maker
}

func NewAuthMiddlewareBuilder(maker token.Maker) *AuthMiddlewareBuilder {
	return &AuthMiddlewareBuilder{
		maker: maker,
	}
}

func (b *AuthMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResp(err))
			return
		}

		fds := strings.Fields(authorizationHeader)
		if len(fds) < 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResp(err))
			return
		}

		if fds[0] != authorizationTyp {
			err := fmt.Errorf("unsupported authorization type: %s", fds[0])
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResp(err))
			return
		}

		payload, err := b.maker.Verify(fds[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResp(err))
			return
		}

		ctx.Set(payloadKey, payload)
		ctx.Next()
	}
}
