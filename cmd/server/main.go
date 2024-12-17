package main

import (
	"context"
	"os"

	"github.com/authgear/authgear-nft-indexer/cmd/server/cmd"
	_ "github.com/authgear/authgear-nft-indexer/cmd/server/cmd/cmddatabase"
	_ "github.com/authgear/authgear-nft-indexer/cmd/server/cmd/cmdstart"
)

func main() {
	ctx := context.Background()
	err := cmd.Root.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
