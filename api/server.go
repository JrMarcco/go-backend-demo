package api

import (
	"github.com/gin-gonic/gin"
	db "go-backend-demo/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(s *db.Store) *Server {
	r := gin.Default()

	server := &Server{
		store:  s,
		router: r,
	}

	r.POST("/account", server.createAccount)
	r.GET("/account/:id", server.getAccount)
	r.GET("/account", server.listAccount)

	return server
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResp(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
