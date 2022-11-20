package api

func (s *S) RegisterRouter() {
	// use for health check
	s.router.GET("/healthz", s.healthz)

	s.router.POST("/user/add", s.createUser)
	s.router.POST("/login", s.login)

	apiG := s.router.Group("/api/v1")
	apiG.Use(NewAuthMiddlewareBuilder(s.tokenMaker).Build())

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
		userG.GET("/get/:id", s.getUser)
	}
}
