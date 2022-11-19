package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) healthz(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
