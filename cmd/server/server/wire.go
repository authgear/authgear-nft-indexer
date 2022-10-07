//go:build wireinject
// +build wireinject

package server

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/google/wire"
)

func NewHealthCheckAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.HealthCheckAPIHandler))))
}

func NewListCollectionAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.ListCollectionAPIHandler))))
}

func NewListOwnerNFTAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.ListOwnerNFTAPIHandler))))
}

func NewGetCollectionAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.GetCollectionAPIHandler))))
}

func NewGetCollectionMetadataAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.GetCollectionMetadataAPIHandler))))
}
