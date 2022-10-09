package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/authgear/authgear-nft-indexer/pkg/mutator"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/ratelimit"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
	agratelimit "github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/httputil"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	query.DependencySet,
	wire.Bind(new(handler.ListOwnerNFTHandlerNFTOwnerQuery), new(*query.NFTOwnerQuery)),

	mutator.DependencySet,
	wire.Bind(new(handler.ListOwnerNFTHandlerNFTOwnershipMutator), new(*mutator.NFTOwnershipMutator)),
	wire.Bind(new(handler.GetCollectionMetadataNFTCollectionMutator), new(*mutator.NFTCollectionMutator)),

	web3.DependencySet,
	wire.Bind(new(handler.GetCollectionMetadataHandlerAlchemyAPI), new(*web3.AlchemyAPI)),
	wire.Bind(new(handler.ProbeCollectionHandlerAlchemyAPI), new(*web3.AlchemyAPI)),
	wire.Bind(new(handler.ListOwnerNFTHandlerAlchemyAPI), new(*web3.AlchemyAPI)),

	ratelimit.DependencySet,
	wire.Bind(new(handler.GetCollectionMetadataRateLimiter), new(*agratelimit.Limiter)),
	wire.Bind(new(handler.ProbeCollectionHandlerRateLimiter), new(*agratelimit.Limiter)),

	handler.DependencySet,
	httputil.DependencySet,
	wire.Bind(new(handler.JSONResponseWriter), new(*httputil.JSONResponseWriter)),
)
