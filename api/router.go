package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	router := gin.Default()

	//create account router
	router.POST("/accounts", server.createAccount)
	//get account by id
	router.GET("/accounts/:id", server.getAccount)
	//get all accounts
	router.GET("/accounts", server.listAccounts)
	//update accounts by id
	router.PUT("/accounts/:id", server.updateAccounts)
	router.DELETE("/accounts/:id", server.deleteAccounts)

	// transfer tx from one account to another
	router.POST("/transfer", server.createTransfer)

	// Create user
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}
