//go:build wireinject
// +build wireinject

package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/google/wire"
)

func NewServer(
	config config.Config,
) Server {
	panic(wire.Build(DependencySet))
}
