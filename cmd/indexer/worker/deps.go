package worker

import (
	"github.com/authgear/authgear-nft-indexer/pkg/mutator"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
	"github.com/authgear/authgear-nft-indexer/pkg/worker/task"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	query.DependencySet,
	wire.Bind(new(task.SycnNFTCollectionTaskCollectionQuery), new(*query.NFTCollectionQuery)),

	mutator.DependencySet,
	wire.Bind(new(task.SycnNFTTransferTransferMutator), new(*mutator.NFTTransferMutator)),

	web3.DependencySet,
	wire.Bind(new(task.SycnNFTTransferAlchemyAPI), new(*web3.AlchemyAPI)),

	task.DependencySet,
	wire.Struct(new(Worker), "*"),
)
