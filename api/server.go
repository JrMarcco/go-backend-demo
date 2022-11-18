package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"github.com/jrmarcco/go-backend-demo/token"
	"github.com/jrmarcco/go-backend-demo/util"
)

type Server struct {
	router *gin.Engine

	config util.ServerCfg
	store  db.Store

	tokenMaker token.Maker
}

func NewServer(config util.ServerCfg, s db.Store) *Server {
	r := gin.Default()

	server := &Server{
		router:     r,
		config:     config,
		store:      s,
		tokenMaker: token.NewPasetoPubMarkerV4(),
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

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	_ = s.router.Use(middleware...)
}

func (s *Server) GenerateToken(username string) (string, error) {
	return s.tokenMaker.Generate(username, s.config.TokenDuration)
}

func (s *Server) VerifyToken(token string) (*token.Payload, error) {
	return s.tokenMaker.Verify(token)
}

func ErrorResp(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
