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

func NewListOwnerNFTAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.ListOwnerNFTAPIHandler))))
}

func NewGetCollectionMetadataAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.GetCollectionMetadataAPIHandler))))
}

func NewProbeCollectionAPIHandler(
	p *handler.RequestProvider,
) http.Handler {
	panic(wire.Build(DependencySet, wire.Bind(new(http.Handler), new(*handler.ProbeCollectionAPIHandler))))
}
