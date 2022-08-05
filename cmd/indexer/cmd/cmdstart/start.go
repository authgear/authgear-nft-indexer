package cmdstart

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/authgear/authgear-nft-indexer/cmd/indexer/cmd"
	indexercmd "github.com/authgear/authgear-nft-indexer/cmd/indexer/cmd"
	"github.com/authgear/authgear-nft-indexer/cmd/indexer/worker"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
)

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "Start indexer",
	RunE: func(cmd *cobra.Command, args []string) error {
		binder := indexercmd.GetBinder()
		configPath, err := binder.GetRequiredString(cmd, indexercmd.ArgConfig)
		if err != nil {
			return err
		}

		ctx := context.Background()
		config := config.NewConfig(configPath)

		worker := worker.NewWorker(ctx, config)

		worker.Start()
		return nil
	},
}

func init() {
	binder := indexercmd.GetBinder()
	binder.BindString(cmdStart.Flags(), indexercmd.ArgConfig)
	cmd.Root.AddCommand(cmdStart)
}
