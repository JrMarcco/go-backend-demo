package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(s db.Store) *Server {
	r := gin.Default()

	server := &Server{
		store:  s,
		router: r,
	}

	return server
}

func (s *Server) Start(addr string) error {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			return err
		}
	}

	return s.router.Run(addr)
}

func errorResp(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
