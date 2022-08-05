//go:build wireinject
// +build wireinject

package worker

import (
	"context"
	"github.com/google/wire"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-nft-indexer/pkg/worker/task"
	"github.com/uptrace/bun"
)

func NewWorker(
	ctx context.Context,
	config config.Config,
) Worker {
	panic(wire.Build(DependencySet))
}

func NewSyncETHNFTCollectionTaskHandler(
	ctx context.Context,
	config config.Config,
	session *bun.DB,
) model.Task {
	panic(wire.Build(
		DependencySet,
		wire.Bind(new(model.Task), new(*task.SyncETHNFTCollectionTaskHandler)),
	))
}

func NewSyncETHNFTTransferTaskHandler(
	ctx context.Context,
	config config.Config,
	session *bun.DB,
) model.Task {
	panic(wire.Build(
		DependencySet,
		wire.Bind(new(model.Task), new(*task.SyncETHNFTTransferTaskHandler)),
	))
}
