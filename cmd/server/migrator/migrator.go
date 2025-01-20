package migrator

import (
	"embed"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-server/pkg/util/sqlmigrate"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	Config config.Config
}

//go:embed migrations
var mainMigrationSetFS embed.FS

var MainMigrationSet = sqlmigrate.NewMigrateSet(sqlmigrate.NewMigrationSetOptions{
	TableName:                            "migration",
	EmbedFS:                              mainMigrationSetFS,
	EmbedFSRoot:                          "migrations",
	OutputPathRelativeToWorkingDirectory: "./cmd/server/migrator/migrations",
})

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
