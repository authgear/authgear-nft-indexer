package ratelimit

import (
	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet(
	ratelimit.NewLogger,
	clock.DependencySet,

	wire.Struct(new(Factory), "*"),
)
