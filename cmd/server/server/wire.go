//go:build wireinject
// +build wireinject

package server

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/google/wire"
)

func NewRegisterCollectionAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.RegisterCollectionAPIHandler))))
}

func NewDeregisterCollectionAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.DeregisterCollectionAPIHandler))))
}

func NewListCollectionAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.ListCollectionAPIHandler))))
}

func NewListCollectionOwnerAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.ListCollectionOwnersAPIHandler))))
}

func NewListOwnerNFTAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.ListOwnerNFTAPIHandler))))
}
