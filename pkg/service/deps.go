package service

import (
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	wire.Struct(new(MetadataService), "*"),
	wire.Struct(new(ProbeService), "*"),
	wire.Struct(new(OwnershipService), "*"),
)
