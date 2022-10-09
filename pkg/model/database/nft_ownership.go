package database

import (
	"time"

	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
)

type NFTOwnership struct {
	bun.BaseModel `bun:"table:eth_nft_ownership"`
	Base

	Blockchain      string      `bun:"blockchain,notnull"`
	Network         string      `bun:"network,notnull"`
	ContractAddress string      `bun:"contract_address,notnull"`
	TokenID         string      `bun:"token_id,notnull"`
	Balance         string      `bun:"balance,notnull"`
	BlockNumber     *bunbig.Int `bun:"block_number,notnull"`
	OwnerAddress    string      `bun:"owner_address,notnull"`
	TransactionHash string      `bun:"txn_hash,notnull"`
	BlockTimestamp  time.Time   `bun:"block_timestamp,notnull"`
}

func (c NFTOwnership) ContractID() authgearweb3.ContractID {
	return authgearweb3.ContractID{
		Blockchain:      c.Blockchain,
		Network:         c.Network,
		ContractAddress: c.ContractAddress,
	}
}
