package handler

import (
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	wire.Struct(new(RegisterCollectionAPIHandler), "*"),
	wire.Struct(new(DeregisterCollectionAPIHandler), "*"),
	wire.Struct(new(ListCollectionAPIHandler), "*"),
	wire.Struct(new(ListCollectionOwnersAPIHandler), "*"),
)
