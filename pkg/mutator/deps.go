package mutator

import (
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	wire.Struct(new(NFTCollectionMutator), "*"),
	wire.Struct(new(NFTOwnershipMutator), "*"),
	wire.Struct(new(NFTCollectionProbeMutator), "*"),
)
