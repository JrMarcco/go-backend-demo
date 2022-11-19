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
	Router     *gin.Engine
	TokenMaker token.Maker

	config util.ServerCfg
	store  db.Store
}

func NewServer(config util.ServerCfg, s db.Store) *Server {
	r := gin.Default()

	server := &Server{
		Router:     r,
		config:     config,
		store:      s,
		TokenMaker: token.NewPasetoPubMarkerV4(),
	}

	return server
}

func (s *Server) Start(addr string) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			return err
		}
	}

	return s.Router.Run(addr)
}

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	_ = s.Router.Use(middleware...)
}

func (s *Server) GenerateToken(username string) (string, error) {
	return s.TokenMaker.Generate(username, s.config.TokenDuration)
}

func (s *Server) VerifyToken(token string) (*token.Payload, error) {
	return s.TokenMaker.Verify(token)
}

func ErrorResp(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
