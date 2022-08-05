package main

import (
	"os"

	"github.com/authgear/authgear-nft-indexer/cmd/indexer/cmd"
	_ "github.com/authgear/authgear-nft-indexer/cmd/indexer/cmd/cmddatabase"
	_ "github.com/authgear/authgear-nft-indexer/cmd/indexer/cmd/cmdstart"
)

func main() {
	err := cmd.Root.Execute()
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
