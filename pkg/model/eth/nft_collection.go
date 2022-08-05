package eth

import (
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"

	"github.com/authgear/authgear-nft-indexer/pkg/model"
)

type NFTCollection struct {
	bun.BaseModel `bun:"table:eth_nft_collection"`
	model.BaseWithID

	Blockchain        string     `bun:"blockchain,notnull"`
	Network           string     `bun:"network,notnull"`
	ContractAddress   string     `bun:"contract_address,notnull"`
	Name              string     `bun:"name,notnull"`
	SyncedBlockHeight bunbig.Int `bun:"synced_block_height,notnull"`
}
