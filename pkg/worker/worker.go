package worker

import (
	"strconv"
	"strings"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/jrallison/go-workers"
)

func ConfigureWorkers(config config.RedisConfig) error {
	url, err := config.RedisURL()
	if err != nil {
		return err
	}

	options := map[string]string{
		// location of redis instance
		"server": url.Host,

		// instance of the database
		"database": strings.TrimPrefix(url.Path, "/"),

		// number of connections to keep open with redis
		"pool": strconv.Itoa(config.MaxOpenConnection),
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "nft-indexer",
	}

	if url.User != nil {
		if p, ok := url.User.Password(); ok {
			options["password"] = p
		}
	}

	workers.Configure(options)

	return nil
}
