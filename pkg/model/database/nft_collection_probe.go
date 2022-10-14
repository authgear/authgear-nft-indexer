package database

import (
	"github.com/uptrace/bun"

	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

type NFTCollectionProbe struct {
	bun.BaseModel `bun:"table:eth_nft_collection_probe"`

	Blockchain        string             `bun:"blockchain,notnull"`
	Network           string             `bun:"network,notnull"`
	ContractAddress   authgearweb3.EIP55 `bun:"contract_address,notnull"`
	IsLargeCollection bool               `bun:"is_large_collection,notnull"`
}
