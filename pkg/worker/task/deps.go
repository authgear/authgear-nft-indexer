package task

import (
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	wire.Struct(new(SyncETHNFTCollectionTaskHandler), "*"),
	wire.Struct(new(SyncETHNFTTransferTaskHandler), "*"),
)
