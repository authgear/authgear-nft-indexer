//go:build wireinject
// +build wireinject

package migrator

import (
	"github.com/google/wire"
)

func NewMigrator(
	configPath string,
) Migrator {
	panic(wire.Build(DependencySet))
}
