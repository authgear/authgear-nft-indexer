package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/database"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/authgear/authgear-server/pkg/util/server"
)

type Controller struct {
	Config config.Config
	logger *log.Logger
}

func (c *Controller) Start() {
	u, err := server.ParseListenAddress(c.Config.Server.ListenAddr)
	if err != nil {
		c.logger.WithError(err).Fatal("failed to parse admin API server listen address")
	}

	database := database.GetDatabase(c.Config.Database)

	lf := log.NewFactory(log.LevelInfo)
	c.logger = lf.New("server")

	server.Start(c.logger, []server.Spec{
		{
			Name:          "Indexer API Server",
			ListenAddress: u.Host,
			Handler: NewRouter(
				c.Config,
				database,
				lf,
			),
		},
	})
}
