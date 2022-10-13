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

func (c NFTOwnership) ContractTokenID() (*authgearweb3.ContractID, error) {
	values := url.Values{}
	values.Set("token_ids", c.TokenID)
	return authgearweb3.NewContractID(c.Blockchain, c.Network, c.ContractAddress, values)
}

func (c NFTOwnership) ToAPIToken() apimodel.Token {
	return apimodel.Token{
		TokenID: c.TokenID,
		TransactionIdentifier: apimodel.TransactionIdentifier{
			Hash: c.TransactionHash,
		},
		BlockIdentifier: apimodel.BlockIdentifier{
			Index:     *c.BlockNumber.ToMathBig(),
			Timestamp: c.BlockTimestamp,
		},
		Balance: c.Balance,
	}
}

func (c NFTOwnership) IsEmpty() bool {
	return c.Balance == "0" || c.TokenID == "0x0" || c.BlockNumber == bunbig.FromInt64(0) || c.TransactionHash == "0x0" || c.TransactionIndex == 0 || c.BlockTimestamp == nil
}

func NewEmptyNFTOwnership(contractID authgearweb3.ContractID, tokenID string, ownerID authgearweb3.ContractID) NFTOwnership {
	return NFTOwnership{
		Blockchain:       contractID.Blockchain,
		Network:          contractID.Network,
		ContractAddress:  contractID.Address,
		OwnerAddress:     ownerID.Address,
		TokenID:          tokenID,
		Balance:          "0",
		BlockNumber:      bunbig.FromInt64(0),
		TransactionHash:  "0x0",
		TransactionIndex: 0,
		BlockTimestamp:   nil,
	}
}
