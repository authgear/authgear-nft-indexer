package cmd

import (
	"github.com/authgear/authgear-server/pkg/util/cobraviper"
)

var cvbinder *cobraviper.Binder

func GetBinder() *cobraviper.Binder {
	if cvbinder == nil {
		cvbinder = cobraviper.NewBinder()
	}
	return cvbinder
}

var ArgConfig = &cobraviper.StringArgument{
	ArgumentName: "config",
	EnvName:      "CONFIG",
	Usage:        "Config Path",
	DefaultValue: "authgear-nft-indexer.yaml",
}
