package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/database"
	"github.com/authgear/authgear-nft-indexer/pkg/worker"
	agconfig "github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/infra/redis"
	"github.com/authgear/authgear-server/pkg/lib/infra/redis/appredis"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/authgear/authgear-server/pkg/util/server"
)

type Controller struct {
	Config config.Config
	logger *log.Logger
}

func (c *Controller) Start() {
	worker.ConfigureWorkers(c.Config.Redis)

	u, err := server.ParseListenAddress(c.Config.Server.ListenAddr)
	if err != nil {
		c.logger.WithError(err).Fatal("failed to parse admin API server listen address")
	}

	database := database.GetDatabase(c.Config.Database)

	lf := log.NewFactory(log.LevelInfo)
	c.logger = lf.New("server")

	redisPool := redis.NewPool()
	redisHub := redis.NewHub(redisPool, lf)

	redis := appredis.NewHandle(
		redisPool,
		redisHub,
		&agconfig.RedisEnvironmentConfig{
			MaxOpenConnection:     c.Config.Redis.MaxOpenConnection,
			MaxIdleConnection:     c.Config.Redis.MaxIdleConnection,
			MaxConnectionLifetime: agconfig.DurationSeconds(c.Config.Redis.MaxConnectionLifeTime),
			IdleConnectionTimeout: agconfig.DurationSeconds(c.Config.Redis.IdleConnectionTimeout),
		},
		&agconfig.RedisCredentials{
			RedisURL: c.Config.Redis.RedisURL(),
		},
		lf,
	)

	server.Start(c.logger, []server.Spec{
		{
			Name:          "Indexer API Server",
			ListenAddress: u.Host,
			Handler: NewRouter(
				c.Config,
				database,
				redis,
				lf,
			),
		},
	})
}
