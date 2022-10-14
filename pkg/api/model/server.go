package model

import (
	"math/big"
	"time"

	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

type NFTCollection struct {
	ID              string             `json:"id"`
	Blockchain      string             `json:"blockchain"`
	Network         string             `json:"network"`
	Name            string             `json:"name"`
	ContractAddress authgearweb3.EIP55 `json:"contract_address"`
	TotalSupply     *big.Int           `json:"total_supply"`
	Type            string             `json:"type"`
}

type AccountIdentifier struct {
	Address authgearweb3.EIP55 `json:"address"`
}

type NetworkIdentifier struct {
	Blockchain string `json:"blockchain"`
	Network    string `json:"network"`
}

type Contract struct {
	Name    string             `json:"name"`
	Address authgearweb3.EIP55 `json:"address"`
	Type    string             `json:"type"`
}

type TransactionIdentifier struct {
	Hash string `json:"hash"`
}

type BlockIdentifier struct {
	Index     big.Int    `json:"index"`
	Timestamp *time.Time `json:"timestamp"`
}

type Token struct {
	TokenID               string                `json:"token_id"`
	TransactionIdentifier TransactionIdentifier `json:"transaction_identifier"`
	BlockIdentifier       BlockIdentifier       `json:"block_identifier"`
	Balance               string                `json:"balance"`
}

type NFT struct {
	Contract Contract `json:"contract"`
	Tokens   []Token  `json:"tokens"`
}
type NFTOwnership struct {
	AccountIdentifier AccountIdentifier `json:"account_identifier"`
	NetworkIdentifier NetworkIdentifier `json:"network_identifier"`
	NFTs              []NFT             `json:"nfts"`
}

func NewNFTOwnership(ownerID authgearweb3.ContractID, nfts []NFT) NFTOwnership {
	return NFTOwnership{
		AccountIdentifier: AccountIdentifier{
			Address: ownerID.Address,
		},
		NetworkIdentifier: NetworkIdentifier{
			Blockchain: ownerID.Blockchain,
			Network:    ownerID.Network,
		},
		NFTs: nfts,
	}
}

type GetContractMetadataRequestData struct {
	ContractIDs []authgearweb3.ContractID `json:"contract_ids"`
}
type GetContractMetadataResponse struct {
	Collections []NFTCollection `json:"collections"`
}

type ProbeCollectionRequestData struct {
	ContractID authgearweb3.ContractID `json:"contract_id"`
}

type ProbeCollectionResponse struct {
	IsLargeCollection bool `json:"is_large_collection"`
}

type ListOwnerNFTRequestData struct {
	OwnerAddress authgearweb3.ContractID   `json:"owner_address"`
	ContractIDs  []authgearweb3.ContractID `json:"contract_ids"`
}
