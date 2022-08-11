// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/config"
)

// Injectors from wire.go:

func NewServer(config2 config.Config) Server {
	server := Server{
		config: config2,
	}
	return server
}