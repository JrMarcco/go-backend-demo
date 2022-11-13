package api

import (
	"github.com/gin-gonic/gin"
	"go-backend-demo/db"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}
