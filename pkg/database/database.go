package database

import (
	"database/sql"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func GetDatabase(config config.DatabaseConfig) *bun.DB {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(config.URL),
	))

	sqlDb.SetMaxOpenConns(config.PoolSize)

	db := bun.NewDB(sqlDb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	if config.Verbose {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db
}
