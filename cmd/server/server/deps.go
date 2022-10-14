package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/authgear/authgear-nft-indexer/pkg/mutator"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/ratelimit"
	"github.com/authgear/authgear-nft-indexer/pkg/service"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
	agratelimit "github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/httputil"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(

	query.DependencySet,
	wire.Bind(new(service.ProbeServiceNFTCollectionProbeQuery), new(*query.NFTCollectionProbeQuery)),

	mutator.DependencySet,
	wire.Bind(new(service.MetadataServiceNFTCollectionMutator), new(*mutator.NFTCollectionMutator)),
	wire.Bind(new(service.ProbeServiceNFTCollectionProbeMutator), new(*mutator.NFTCollectionProbeMutator)),
	wire.Bind(new(service.OwnershipServiceNFTOwnershipMutator), new(*mutator.NFTOwnershipMutator)),

	web3.DependencySet,
	wire.Bind(new(service.MetadataServiceAlchemyAPI), new(*web3.AlchemyAPI)),
	wire.Bind(new(service.ProbeServiceAlchemyAPI), new(*web3.AlchemyAPI)),
	wire.Bind(new(service.OwnershipServiceAlchemyAPI), new(*web3.AlchemyAPI)),

	ratelimit.DependencySet,
	wire.Bind(new(handler.GetCollectionMetadataHandlerRateLimiter), new(*agratelimit.Limiter)),
	wire.Bind(new(handler.ProbeCollectionHandlerRateLimiter), new(*agratelimit.Limiter)),

	service.DependencySet,
	wire.Bind(new(handler.GetCollectionMetadataHandlerMetadataService), new(*service.MetadataService)),
	wire.Bind(new(handler.ListOwnerNFTHandlerMetadataService), new(*service.MetadataService)),
	wire.Bind(new(handler.ProbeCollectionHandlerProbeService), new(*service.ProbeService)),
	wire.Bind(new(handler.ListOwnerNFTHandlerOwnershipService), new(*service.OwnershipService)),

	handler.DependencySet,
	httputil.DependencySet,
	wire.Bind(new(handler.JSONResponseWriter), new(*httputil.JSONResponseWriter)),
)
