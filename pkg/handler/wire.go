//go:build wireinject
// +build wireinject

package handler

import (
	"github.com/authgear/authgear-nft-indexer/pkg/ratelimit"
	"github.com/authgear/authgear-server/pkg/lib/infra/redis/appredis"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/google/wire"
)

func NewRateLimiterFactory(
	lf *log.Factory,
	redis *appredis.Handle,
) ratelimit.Factory {
	panic(wire.Build(ratelimit.DependencySet))
}
