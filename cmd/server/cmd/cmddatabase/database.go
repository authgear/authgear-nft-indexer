package cmddatabase

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	servercmd "github.com/authgear/authgear-nft-indexer/cmd/server/cmd"
	"github.com/authgear/authgear-nft-indexer/cmd/server/migrator"
)

func init() {
	binder := servercmd.GetBinder()
	cmdDatabase.AddCommand(cmdMigrate)

	cmdMigrate.AddCommand(cmdMigrateNew)
	cmdMigrate.AddCommand(cmdMigrateUp)
	cmdMigrate.AddCommand(cmdMigrateDown)
	cmdMigrate.AddCommand(cmdMigrateStatus)

	for _, cmd := range []*cobra.Command{cmdMigrateNew, cmdMigrateUp, cmdMigrateDown, cmdMigrateStatus} {
		binder.BindString(cmd.Flags(), servercmd.ArgConfig)
	}

	servercmd.Root.AddCommand(cmdDatabase)
}

var cmdDatabase = &cobra.Command{
	Use:   "database migrate",
	Short: "Database commands",
}

var cmdMigrate = &cobra.Command{
	Use:   "migrate [new|status|up|down]",
	Short: "Migrate database schema",
}

var cmdMigrateNew = &cobra.Command{
	Use:    "new",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		binder := servercmd.GetBinder()
		configPath, err := binder.GetRequiredString(cmd, servercmd.ArgConfig)
		if err != nil {
			return err
		}

		name := strings.Join(args, "_")
		migrator := migrator.NewMigrator(configPath)

		_, err = migrator.Create(name)
		if err != nil {
			return
		}

		return
	},
}

var cmdMigrateUp = &cobra.Command{
	Use:   "up",
	Short: "Migrate database schema to latest version",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		binder := servercmd.GetBinder()
		configPath, err := binder.GetRequiredString(cmd, servercmd.ArgConfig)
		if err != nil {
			return err
		}

		migrator := migrator.NewMigrator(configPath)
		_, err = migrator.Up()
		if err != nil {
			return err
		}

		return
	},
}

var cmdMigrateDown = &cobra.Command{
	Use:    "down",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		binder := servercmd.GetBinder()
		configPath, err := binder.GetRequiredString(cmd, servercmd.ArgConfig)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			err = fmt.Errorf("number of migrations to revert not specified; specify 'all' to revert all migrations")
			return
		}

		var numMigrations int
		if args[0] == "all" {
			numMigrations = 0
		} else {
			numMigrations, err = strconv.Atoi(args[0])
			if err != nil {
				err = fmt.Errorf("invalid number of migrations specified: %w", err)
				return
			} else if numMigrations <= 0 {
				err = fmt.Errorf("no migrations specified to revert")
				return
			}
		}

		migrator := migrator.NewMigrator(configPath)
		_, err = migrator.Down(numMigrations)
		if err != nil {
			return
		}

		return
	},
}

var cmdMigrateStatus = &cobra.Command{
	Use:   "status",
	Short: "Get database schema migration status",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		binder := servercmd.GetBinder()
		configPath, err := binder.GetRequiredString(cmd, servercmd.ArgConfig)
		if err != nil {
			return err
		}

		migrator := migrator.NewMigrator(configPath)
		plans, err := migrator.Status()
		if err != nil {
			return
		}

		if len(plans) != 0 {
			os.Exit(1)
		}

		return
	},
}
