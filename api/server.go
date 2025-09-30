package api

import (
	"fmt"
	db "github.com/QuangPham789/simple-bank/db/sqlc"
	"github.com/QuangPham789/simple-bank/token"
	"github.com/QuangPham789/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Server struct {
	store  db.Store
	config util.Config
	token  token.Maker
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{store: store, config: config, token: tokenMaker}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrencies)
	}
	// add router
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)
	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.token))
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.POST("/transfer", server.createTransfer)
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	server.router = router
}
