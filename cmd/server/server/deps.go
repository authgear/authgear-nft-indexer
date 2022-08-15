package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/authgear/authgear-nft-indexer/pkg/mutator"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
	"github.com/authgear/authgear-server/pkg/util/httputil"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/google/wire"
)

func NewLoggerFactory() *log.Factory {
	return log.NewFactory(log.LevelInfo)
}

var DependencySet = wire.NewSet(
	NewLoggerFactory,
	query.DependencySet,
	wire.Bind(new(handler.ListCollectionHandlerCollectionsQuery), new(*query.NFTCollectionQuery)),

	mutator.DependencySet,
	wire.Bind(new(handler.DeregisterCollectionHandlerNFTCollectionMutator), new(*mutator.NFTCollectionMutator)),
	wire.Bind(new(handler.RegisterCollectionHandlerNFTCollectionMutator), new(*mutator.NFTCollectionMutator)),

	web3.DependencySet,
	wire.Bind(new(handler.RegisterCollectionHandlerAlchemyAPI), new(*web3.AlchemyAPI)),

	handler.DependencySet,
	httputil.DependencySet,
	wire.Bind(new(handler.JSONResponseWriter), new(*httputil.JSONResponseWriter)),
)
