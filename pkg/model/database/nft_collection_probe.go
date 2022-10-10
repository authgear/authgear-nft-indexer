package database

import (
	"github.com/uptrace/bun"
)

type NFTCollectionProbe struct {
	bun.BaseModel `bun:"table:nft_collection_probe"`

	Blockchain        string `bun:"blockchain,notnull"`
	Network           string `bun:"network,notnull"`
	ContractAddress   string `bun:"contract_address,notnull"`
	IsLargeCollection bool   `bun:"is_large_collection,notnull"`
}
