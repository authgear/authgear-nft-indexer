package server

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/database"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/authgear/authgear-server/pkg/util/server"
	"github.com/authgear/authgear-server/pkg/util/signalutil"
)

type Controller struct {
	Config config.Config
	logger *log.Logger
}

func (c *Controller) Start(ctx context.Context) {
	u, err := server.ParseListenAddress(c.Config.Server.ListenAddr)
	if err != nil {
		c.logger.WithError(err).Fatal("failed to parse admin API server listen address")
	}

	database := database.GetDatabase(c.Config.Database)

	lf := log.NewFactory(log.LevelInfo)
	c.logger = lf.New("server")

	signalutil.Start(ctx, c.logger, []signalutil.Daemon{
		server.NewSpec(ctx, &server.Spec{
			Name:          "Indexer API Server",
			ListenAddress: u.Host,
			Handler: NewRouter(
				c.Config,
				database,
				lf,
			),
		}),
	}...)
}
