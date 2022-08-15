package cmdstart

import (
	"github.com/authgear/authgear-nft-indexer/cmd/server/cmd"
	indexercmd "github.com/authgear/authgear-nft-indexer/cmd/server/cmd"
	"github.com/authgear/authgear-nft-indexer/cmd/server/server"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/spf13/cobra"
)

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "Start API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		binder := indexercmd.GetBinder()
		configPath, err := binder.GetRequiredString(cmd, indexercmd.ArgConfig)
		if err != nil {
			return err
		}
		config := config.NewConfig(configPath)

		ctrl := server.Controller{
			Config: config,
		}

		ctrl.Start()
		return nil
	},
}

func init() {
	binder := indexercmd.GetBinder()
	binder.BindString(cmdStart.Flags(), indexercmd.ArgConfig)
	cmd.Root.AddCommand(cmdStart)
}
