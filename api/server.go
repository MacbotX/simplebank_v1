package api

import (
	db "github.com/MacbotX/simplebank_v1/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP request for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// to register validator with gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

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

	server.router = router
	return server
}

// Start runs the HTTP server on a specification address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	// // Split the error message into individual validation errors
	// lines := strings.Split(err.Error(), "\n")

	// // Prepare a slice of maps to hold individual field errors
	// var details []map[string]string

	// for _, line := range lines {
	// 	// Example line: "Key: 'tansferRequest.FromAccountID' Error:Field validation for 'FromAccountID' failed on the 'required' tag"
	// 	parts := strings.SplitN(line, "Error:", 2)
	// 	if len(parts) == 2 {
	// 		fieldPart := strings.TrimSpace(parts[0])
	// 		errorMsg := strings.TrimSpace(parts[1])

	// 		details = append(details, map[string]string{
	// 			"key":   fieldPart,
	// 			"error": errorMsg,
	// 		})
	// 	}
	// }

	// return gin.H{
	// 	"errors": details,
	// }
	return gin.H{"error": err.Error()}

}
