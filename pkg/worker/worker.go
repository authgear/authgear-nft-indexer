package worker

import (
	"strconv"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/jrallison/go-workers"
)

func ConfigureWorkers(config config.RedisConfig) {
	workers.Configure(map[string]string{
		// location of redis instance
		"server": config.Server,
		// instance of the database
		"database": config.Database,
		// number of connections to keep open with redis
		"pool": strconv.Itoa(config.MaxOpenConnection),
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "nft-indexer",
	})
}
