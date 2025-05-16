package api

import (
	db "github.com/MacbotX/simplebank_v1/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP request for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}


//NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server  {
	server := &Server{store: store}
	router := gin.Default()

	//create account router\
	router.POST("/accounts", server.createAccount)
	//get account by id
	router.GET("/accounts/:id", server.getAccount)
	//get all accounts 
	router.GET("/accounts", server.listAccounts)
	//update accounts by id
	router.PUT("/accounts/:id", server.updateAccounts)
	router.DELETE("/accounts/:id", server.deleteAccounts)


	server.router = router
	return server
}

//Start runs the HTTP server on a specification address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H  {
	return gin.H{"error": err.Error()}
}