package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/natanaelrusli/go-bank/db/sqlc"
	"github.com/natanaelrusli/go-bank/token"
	"github.com/natanaelrusli/go-bank/util"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker("")
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccount)
	router.GET("/accounts/:id", server.getAccount)

	router.POST("/transfers", server.createTransfer)

	router.POST("/users", server.createUser)

	router.POST("/auth", server.login)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
