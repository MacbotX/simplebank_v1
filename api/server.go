package api

import (
	"fmt"

	db "github.com/MacbotX/simplebank_v1/db/sqlc"
	"github.com/MacbotX/simplebank_v1/token"
	"github.com/MacbotX/simplebank_v1/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP request for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {

	// TokenMaker is set to use pasetoMaker and can be changed to JWTMaker
	tokenMaker, err := token.NewPasetoMaker(config.TokenSynmetricKey)
	
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// router func imported
	server.setupRouter()

	// to register validator with gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	return server, nil
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
