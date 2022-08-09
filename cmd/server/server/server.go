package server

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
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

	router.POST("/register", func(_ *gin.Context) {
		// TODO: Handle collection registartion
	})
	router.POST("/deregister", func(_ *gin.Context) {
		// TODO: Handle collection deregistration
	})

	router.GET("/collections", func(_ *gin.Context) {
		// TODO: Handle listing collections
	})
	router.GET("/collections/:address/owners", func(_ *gin.Context) {
		// TODO: Handle listing owners under a collection
	})
	router.GET("/nfts/:owner", func(_ *gin.Context) {
		// TODO: Handle listing nfts under an address
	})

	panic(router.Run(s.config.Server.ListenAddr))
}
