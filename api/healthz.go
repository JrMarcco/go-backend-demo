package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *S) healthz(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
