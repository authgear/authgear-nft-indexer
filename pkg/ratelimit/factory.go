package ratelimit

import (
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/infra/redis/appredis"
	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/clock"
)

type Factory struct {
	Logger ratelimit.Logger
	Redis  *appredis.Handle
	Clock  clock.Clock
}

func (rlf *Factory) New(appID string) *ratelimit.Limiter {
	sr := &ratelimit.StorageRedis{
		AppID: config.AppID(appID),
		Redis: rlf.Redis,
	}

	return &ratelimit.Limiter{
		Logger:  rlf.Logger,
		Storage: sr,
		Clock:   rlf.Clock,
	}

}
