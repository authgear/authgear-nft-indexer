package migrator

import (
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-server/pkg/util/sqlmigrate"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	Config config.Config
}

var MainMigrationSet = sqlmigrate.NewMigrateSet("migration", "migrations")

func (m *Migrator) Create(name string) (string, error) {
	return MainMigrationSet.Create(name)
}

func (m *Migrator) Up() (int, error) {
	return MainMigrationSet.Up(sqlmigrate.ConnectionOptions{
		DatabaseURL:    m.Config.Database.URL,
		DatabaseSchema: "public",
	}, 0)
}

func (m *Migrator) Down(numMigrations int) (int, error) {
	return MainMigrationSet.Down(sqlmigrate.ConnectionOptions{
		DatabaseURL:    m.Config.Database.URL,
		DatabaseSchema: "public",
	}, numMigrations)
}

func (m *Migrator) Status() ([]*migrate.PlannedMigration, error) {
	return MainMigrationSet.Status(sqlmigrate.ConnectionOptions{
		DatabaseURL:    m.Config.Database.URL,
		DatabaseSchema: "public",
	})
}
