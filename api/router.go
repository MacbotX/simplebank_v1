package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	router := gin.Default()

	// to add some route to make use of the auth middleware
	authRoutes := router.Group("/").Use(authMiddleWare(server.tokenMaker))

	// Accounts router
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PUT("/accounts/:id", server.updateAccounts)
	authRoutes.DELETE("/accounts/:id", server.deleteAccounts)

	// Transfer route
	authRoutes.POST("/transfer", server.createTransfer)

	// User
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}
