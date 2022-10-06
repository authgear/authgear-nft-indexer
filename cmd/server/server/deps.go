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

	mutator.DependencySet,

	web3.DependencySet,
	wire.Bind(new(handler.GetCollectionMetadataHandlerAlchemyAPI), new(*web3.AlchemyAPI)),

	ratelimit.DependencySet,
	wire.Bind(new(handler.GetCollectionMeatadataRateLimiter), new(*agratelimit.Limiter)),

	handler.DependencySet,
	httputil.DependencySet,
	wire.Bind(new(handler.JSONResponseWriter), new(*httputil.JSONResponseWriter)),
)
