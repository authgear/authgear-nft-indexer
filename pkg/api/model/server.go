package model

import (
	"math/big"
	"time"
)

type NFTCollection struct {
	ID              string    `json:"id"`
	Blockchain      string    `json:"blockchain"`
	Network         string    `json:"network"`
	Name            string    `json:"name"`
	ContractAddress string    `json:"contract_address"`
	TotalSupply     *big.Int  `json:"total_supply"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
}

type WatchCollectionRequestData struct {
	ContractID string `json:"contract_id"`
	Name       string `json:"name,omitempty"`
}

type CollectionListResponse struct {
	Items []NFTCollection `json:"items"`
}

type AccountIdentifier struct {
	Address string `json:"address"`
}

type NetworkIdentifier struct {
	Blockchain string `json:"blockchain"`
	Network    string `json:"network"`
}

type Contract struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type TransactionIdentifier struct {
	Hash string `json:"hash"`
}

type BlockIdentifier struct {
	Index     big.Int   `json:"index"`
	Timestamp time.Time `json:"timestamp"`
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

type GetContractMetadataResponse struct {
	Collections []NFTCollection `json:"collections"`
}

type ProbeCollectionResponse struct {
	IsLargeCollection bool `json:"is_large_collection"`
}

type ProbeCollectionRequestData struct {
	AppID      string `json:"app_id"`
	ContractID string `json:"contract_id"`
}
