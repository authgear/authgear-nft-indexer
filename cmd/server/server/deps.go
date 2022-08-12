package server

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/authgear/authgear-nft-indexer/pkg/mutator"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func ProvideContext(context *gin.Context) context.Context {
	return context.Request.Context()
}

var DependencySet = wire.NewSet(

	query.DependencySet,
	wire.Bind(new(handler.ListCollectionHandlerCollectionsQuery), new(*query.NFTCollectionQuery)),

	mutator.DependencySet,
	wire.Bind(new(handler.DeregisterCollectionHandlerNFTCollectionMutator), new(*mutator.NFTCollectionMutator)),
	wire.Bind(new(handler.RegisterCollectionHandlerNFTCollectionMutator), new(*mutator.NFTCollectionMutator)),

	web3.DependencySet,
	wire.Bind(new(handler.RegisterCollectionHandlerAlchemyAPI), new(*web3.AlchemyAPI)),

	handler.DependencySet,
	ProvideContext,
	wire.Struct(new(Server), "*"),
)
