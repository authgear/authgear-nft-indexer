package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/authgear/authgear-nft-indexer/pkg/mutator"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
	"github.com/authgear/authgear-server/pkg/util/httputil"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	query.DependencySet,
	wire.Bind(new(handler.ListOwnerNFTHandlerNFTCollectionQuery), new(*query.NFTCollectionQuery)),

	mutator.DependencySet,
	wire.Bind(new(handler.WatchCollectionHandlerNFTCollectionMutator), new(*mutator.NFTCollectionMutator)),

	web3.DependencySet,
	wire.Bind(new(handler.WatchCollectionHandlerAlchemyAPI), new(*web3.AlchemyAPI)),

	handler.DependencySet,
	httputil.DependencySet,
	wire.Bind(new(handler.JSONResponseWriter), new(*httputil.JSONResponseWriter)),
)
