package server

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
}

func (s *Server) Start() {
	router := gin.Default()

	// db := database.GetDatabase(s.config.Database)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	registerHandler := handler.NewRegisterCollectionAPIHandler(s.session, s.config)
	router.POST("/register", registerHandler.Handler())

	deregisterHandler := handler.NewDeregisterCollectionAPIHandler(s.session, s.config)
	router.POST("/deregister", deregisterHandler.Handler())

	listHandler := handler.NewListCollectionAPIHandler(s.session, s.config)
	router.GET("/collections", listHandler.Handler())

	router.GET("/collections/:address/owners", func(_ *gin.Context) {
		// TODO: Handle listing owners under a collection
	})
	router.GET("/nfts/:owner", func(_ *gin.Context) {
		// TODO: Handle listing nfts under an address
	})

	panic(router.Run(s.config.Server.ListenAddr))
}
