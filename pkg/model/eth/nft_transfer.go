package eth

import (
	"time"

	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
)

type NFTTransfer struct {
	bun.BaseModel `bun:"table:eth_nft_transfer"`
	model.Base

	Blockchain      string      `bun:"blockchain,notnull"`
	Network         string      `bun:"network,notnull"`
	ContractAddress string      `bun:"contract_address,notnull"`
	TokenID         *bunbig.Int `bun:"token_id,notnull"`
	BlockNumber     *bunbig.Int `bun:"block_number,notnull"`
	FromAddress     string      `bun:"from_address,notnull"`
	ToAddress       string      `bun:"to_address,notnull"`
	TransactionHash string      `bun:"txn_hash,notnull"`
	BlockTimestamp  time.Time   `bun:"block_timestamp,notnull"`
}
