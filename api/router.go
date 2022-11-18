package api

func (s *Server) RegisterRouter() {
	apiG := s.router.Group("/api/v1")

	accountG := apiG.Group("/account")
	{
		accountG.POST("/add", s.createAccount)
		accountG.GET("/get/:id", s.getAccount)
		accountG.GET("/get", s.listAccount)
	}

	transferG := apiG.Group("/transfer")
	{
		transferG.POST("/add", s.createTransfer)
	}

	userG := apiG.Group("/user")
	{
		userG.POST("login", s.login)
		userG.POST("/add", s.createUser)
		userG.GET("/get/:id", s.getUser)
	}
}
