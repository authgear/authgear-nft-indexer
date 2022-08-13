package server

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/database"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
}

func (s *Server) Start() {
	router := gin.Default()

	db := database.GetDatabase(s.config.Database)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.POST("/register", func(ctx *gin.Context) {
		registerHandler := NewRegisterCollectionAPIHandler(ctx, s.config, db)
		registerHandler.Handle()
	})

	router.POST("/deregister", func(ctx *gin.Context) {
		deregisterHandler := NewDeregisterCollectionAPIHandler(ctx, s.config, db)
		deregisterHandler.Handle()
	})

	router.GET("/collections", func(ctx *gin.Context) {
		listHandler := NewListCollectionAPIHandler(ctx, s.config, db)
		listHandler.Handle()
	})

	router.GET("/collections/:blockchain/:network/owners/:contract_address", func(ctx *gin.Context) {
		listOwnerHandler := NewListCollectionOwnerAPIHandler(ctx, s.config, db)
		listOwnerHandler.Handle()
	})
	router.GET("/nfts/:owner_address", func(ctx *gin.Context) {
		listOwnerNFTHandler := NewListOwnerNFTAPIHandler(ctx, s.config, db)
		listOwnerNFTHandler.Handle()
	})

	panic(router.Run(s.config.Server.ListenAddr))
}
