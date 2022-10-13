package database

import (
	"net/url"
	"time"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
)

type NFTOwnership struct {
	bun.BaseModel `bun:"table:eth_nft_ownership"`
	Base

	Blockchain       string      `bun:"blockchain,notnull"`
	Network          string      `bun:"network,notnull"`
	ContractAddress  string      `bun:"contract_address,notnull"`
	TokenID          string      `bun:"token_id,notnull"`
	Balance          string      `bun:"balance,notnull"`
	BlockNumber      *bunbig.Int `bun:"block_number,notnull"`
	OwnerAddress     string      `bun:"owner_address,notnull"`
	TransactionHash  string      `bun:"txn_hash,notnull"`
	TransactionIndex int         `bun:"txn_index,notnull"`
	BlockTimestamp   *time.Time  `bun:"block_timestamp"`
}

func (c NFTOwnership) ContractID() (*authgearweb3.ContractID, error) {
	return authgearweb3.NewContractID(c.Blockchain, c.Network, c.ContractAddress, url.Values{})
}

func (c NFTOwnership) ToAPIToken() apimodel.Token {
	return apimodel.Token{
		TokenID: c.TokenID,
		TransactionIdentifier: apimodel.TransactionIdentifier{
			Hash:  c.TransactionHash,
			Index: c.TransactionIndex,
		},
		BlockIdentifier: apimodel.BlockIdentifier{
			Index:     *c.BlockNumber.ToMathBig(),
			Timestamp: *c.BlockTimestamp,
		},
		Balance: c.Balance,
	}
}
