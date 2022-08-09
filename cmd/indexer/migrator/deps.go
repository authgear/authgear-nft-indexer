package migrator

import (
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	config.NewConfig,
	wire.Struct(new(Migrator), "*"),
)
