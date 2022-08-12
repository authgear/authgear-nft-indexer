//go:build wireinject
// +build wireinject

package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/uptrace/bun"
)

func NewServer(
	config config.Config,
) Server {
	panic(wire.Build(DependencySet))
}

func NewRegisterCollectionAPIHandler(
	ctx *gin.Context,
	config config.Config,
	session *bun.DB,
) handler.RegisterCollectionAPIHandler {
	panic(wire.Build(DependencySet))
}

func NewDeregisterCollectionAPIHandler(
	ctx *gin.Context,
	config config.Config,
	session *bun.DB,
) handler.DeregisterCollectionAPIHandler {
	panic(wire.Build(DependencySet))
}

func NewListCollectionAPIHandler(
	ctx *gin.Context,
	config config.Config,
	session *bun.DB,
) handler.ListCollectionAPIHandler {
	panic(wire.Build(DependencySet))
}

func NewListCollectionOwnerAPIHandler(
	ctx *gin.Context,
	config config.Config,
	session *bun.DB,
) handler.ListCollectionOwnersAPIHandler {
	panic(wire.Build(DependencySet))
}
