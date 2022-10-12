package alchemy

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	"github.com/authgear/authgear-server/pkg/util/hexstring"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun/extra/bunbig"
)

func NewContractTokenIDWithTokenID(blockchain string, network string, address string, tokenID string) (*authgearweb3.ContractID, error) {

	tokenIDHex, err := hexstring.TrimmedParse(tokenID)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("token_ids", tokenIDHex.String())

	return authgearweb3.NewContractID(blockchain, network, address, query)
}

type RawContract struct {
	Value   string `json:"value"`
	Address string `json:"address"`
	Decimal string `json:"decimal"`
}

type Metadata struct {
	BlockTimestamp string `json:"blockTimestamp"`
}

type ERC1155Metadata struct {
	TokenID string `json:"tokenId"`
	Value   string `json:"value"`
}

type TransactionUniqueID struct {
	TransactionHash  string `json:"txn_hash"`
	TransactionIndex int    `json:"txn_index"`
}

func ParseTransactionUniqueID(s string) (*TransactionUniqueID, error) {
	// "0x000000:log:20"
	splitted := strings.Split(s, ":")
	if len(splitted) != 3 || splitted[1] != "log" {
		return nil, fmt.Errorf("failed to parse transaction unique ID")
	}

	logID, err := strconv.Atoi(splitted[2])
	if err != nil {
		return nil, err
	}

	return &TransactionUniqueID{
		TransactionHash:  splitted[0],
		TransactionIndex: logID,
	}, nil

}

type TokenTransfer struct {
	Category        string             `json:"category"`
	UniqueID        string             `json:"uniqueId"`
	Token           string             `json:"token"`
	BlockNum        string             `json:"blockNum"`
	From            string             `json:"from"`
	To              string             `json:"to"`
	Value           string             `json:"value"`
	ERC721TokenID   *string            `json:"erc721TokenId"`
	ERC1155Metadata *[]ERC1155Metadata `json:"erc1155Metadata"`
	TokenID         string             `json:"tokenId"`
	Asset           string             `json:"asset"`
	Hash            string             `json:"hash"`
	RawContract     RawContract        `json:"rawContract"`
	Metadata        Metadata           `json:"metadata"`
}

type AssetTransferRequestParams struct {
	FromBlock         string   `json:"fromBlock,omitempty"`
	ToBlock           string   `json:"toBlock,omitempty"`
	FromAddress       string   `json:"fromAddress,omitempty"`
	ToAddress         string   `json:"toAddress,omitempty"`
	ContractAddresses []string `json:"contractAddresses,omitempty"`
	Category          []string `json:"category"`
	Order             string   `json:"order"`
	WithMetadata      bool     `json:"withMetadata,omitempty"`
	ExcludeZeroValue  bool     `json:"excludeZeroValue,omitempty"`
	MaxCount          string   `json:"maxCount,omitempty"`
	PageKey           string   `json:"pageKey,omitempty"`
}

type AssetTransferResult struct {
	Transfers []TokenTransfer `json:"transfers"`
	PageKey   string          `json:"pageKey"`
}

type AssetTransferError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AssetTransferResponse struct {
	Result AssetTransferResult `json:"result"`
	Error  *AssetTransferError `json:"error,omitempty"`
}

type ContractMetadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply string `json:"totalSupply,omitempty"`
	TokenType   string `json:"tokenType"`
}

type ContractMetadataResponse struct {
	Address          string           `json:"address"`
	ContractMetadata ContractMetadata `json:"contractMetadata"`
}

type GetOwnersForCollectionResponse struct {
	OwnerAddresses []string `json:"ownerAddresses"`
	PageKey        *string  `json:"pageKey,omitempty"`
}

type OwnedNFTContract struct {
	Address string `json:"address"`
}

type OwnedNFTID struct {
	TokenID string `json:"tokenId"`
}
type OwnedNFT struct {
	Contract OwnedNFTContract `json:"contract"`
	ID       OwnedNFTID       `json:"id"`
	Balance  string           `json:"balance"`
}

type GetNFTsResponse struct {
	OwnedNFTs []OwnedNFT `json:"ownedNfts"`
	PageKey   *string    `json:"pageKey,omitempty"`
}

func MakeNFTOwnerships(blockchain string, network string, transfers []TokenTransfer, ownedNFTs []OwnedNFT) ([]database.NFTOwnership, error) {
	contractTokenIDToBalance := make(map[string]string)
	for _, ownedNFT := range ownedNFTs {
		contractTokenID, err := NewContractTokenIDWithTokenID(blockchain, network, ownedNFT.Contract.Address, ownedNFT.ID.TokenID)
		if err != nil {
			return []database.NFTOwnership{}, err
		}

		contractTokenIDURL, err := contractTokenID.URL()
		if err != nil {
			return []database.NFTOwnership{}, err
		}
		contractTokenIDToBalance[contractTokenIDURL.String()] = ownedNFT.Balance
	}

	ownerships := make([]database.NFTOwnership, 0)
	for _, transfer := range transfers {
		blockNumber, err := hexstring.Parse(transfer.BlockNum)
		if err != nil {
			return []database.NFTOwnership{}, err
		}

		blockTime, err := time.Parse(time.RFC3339, transfer.Metadata.BlockTimestamp)
		if err != nil {
			return []database.NFTOwnership{}, err
		}

		uniqueID, err := ParseTransactionUniqueID(transfer.UniqueID)
		if err != nil {
			return []database.NFTOwnership{}, err
		}

		if transfer.ERC1155Metadata != nil {
			for _, erc1155 := range *transfer.ERC1155Metadata {
				contractTokenID, err := NewContractTokenIDWithTokenID(blockchain, network, transfer.RawContract.Address, erc1155.TokenID)
				if err != nil {
					return []database.NFTOwnership{}, err
				}
				contractTokenIDURL, err := contractTokenID.URL()
				if err != nil {
					return []database.NFTOwnership{}, err
				}

				balance := contractTokenIDToBalance[contractTokenIDURL.String()]
				ownerships = append(ownerships, database.NFTOwnership{
					Blockchain:       contractTokenID.Blockchain,
					Network:          contractTokenID.Network,
					ContractAddress:  contractTokenID.Address,
					TokenID:          erc1155.TokenID,
					Balance:          balance,
					BlockNumber:      bunbig.FromMathBig(blockNumber.ToBigInt()),
					OwnerAddress:     transfer.To,
					TransactionHash:  transfer.Hash,
					TransactionIndex: uniqueID.TransactionIndex,
					BlockTimestamp:   blockTime,
				})
			}
			continue
		}

		// Transfer is ERC-721
		contractTokenID, err := NewContractTokenIDWithTokenID(blockchain, network, transfer.RawContract.Address, transfer.TokenID)
		if err != nil {
			return []database.NFTOwnership{}, err
		}
		contractTokenIDURL, err := contractTokenID.URL()
		if err != nil {
			return []database.NFTOwnership{}, err
		}
		balance := contractTokenIDToBalance[contractTokenIDURL.String()]
		ownerships = append(ownerships, database.NFTOwnership{
			Blockchain:       contractTokenID.Blockchain,
			Network:          contractTokenID.Network,
			ContractAddress:  contractTokenID.Address,
			TokenID:          transfer.TokenID,
			Balance:          balance,
			BlockNumber:      bunbig.FromMathBig(blockNumber.ToBigInt()),
			OwnerAddress:     transfer.To,
			TransactionHash:  transfer.Hash,
			TransactionIndex: uniqueID.TransactionIndex,
			BlockTimestamp:   blockTime,
		})
	}

	return ownerships, nil
}
