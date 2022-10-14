package query

import (
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	wire.Struct(new(NFTCollectionQuery), "*"),
	wire.Struct(new(NFTOwnershipQuery), "*"),
	wire.Struct(new(NFTCollectionProbeQuery), "*"),
)
