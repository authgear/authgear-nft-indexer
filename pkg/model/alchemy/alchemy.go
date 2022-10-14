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

type RawContract struct {
	Value   string             `json:"value"`
	Address authgearweb3.EIP55 `json:"address"`
	Decimal string             `json:"decimal"`
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
	From            authgearweb3.EIP55 `json:"from"`
	To              authgearweb3.EIP55 `json:"to"`
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
	FromBlock         string               `json:"fromBlock,omitempty"`
	ToBlock           string               `json:"toBlock,omitempty"`
	FromAddress       authgearweb3.EIP55   `json:"fromAddress,omitempty"`
	ToAddress         authgearweb3.EIP55   `json:"toAddress,omitempty"`
	ContractAddresses []authgearweb3.EIP55 `json:"contractAddresses,omitempty"`
	Category          []string             `json:"category"`
	Order             string               `json:"order"`
	WithMetadata      bool                 `json:"withMetadata,omitempty"`
	ExcludeZeroValue  bool                 `json:"excludeZeroValue,omitempty"`
	MaxCount          string               `json:"maxCount,omitempty"`
	PageKey           string               `json:"pageKey,omitempty"`
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

func MakeNFTOwnerships(ownerID authgearweb3.ContractID, contracts []authgearweb3.ContractID, transfers []TokenTransfer, ownedNFTs []OwnedNFT) ([]database.NFTOwnership, error) {
	contractIDToTokenIDToBalance := make(map[string]map[string]string)
	for _, ownedNFT := range ownedNFTs {
		contractID, err := authgearweb3.NewContractID(ownerID.Blockchain, ownerID.Network, ownedNFT.Contract.Address, url.Values{})
		if err != nil {
			return []database.NFTOwnership{}, err
		}

		contractURL := contractID.String()

		if _, ok := contractIDToTokenIDToBalance[contractURL]; !ok {
			contractIDToTokenIDToBalance[contractURL] = make(map[string]string)
		}
		contractIDToTokenIDToBalance[contractURL][ownedNFT.ID.TokenID] = ownedNFT.Balance
	}

	contractIDToTokenIDToOwnership := make(map[string]map[string]database.NFTOwnership, 0)
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

		contractID, err := authgearweb3.NewContractID(ownerID.Blockchain, ownerID.Network, transfer.RawContract.Address.String(), url.Values{})
		if err != nil {
			return []database.NFTOwnership{}, err
		}

		contractURL := contractID.String()

		if _, ok := contractIDToTokenIDToOwnership[contractURL]; !ok {
			contractIDToTokenIDToOwnership[contractURL] = make(map[string]database.NFTOwnership)
		}

		if transfer.ERC1155Metadata != nil {
			for _, erc1155 := range *transfer.ERC1155Metadata {
				balance := contractIDToTokenIDToBalance[contractURL][erc1155.TokenID]
				if _, ok := contractIDToTokenIDToOwnership[contractURL][erc1155.TokenID]; !ok {

					contractIDToTokenIDToOwnership[contractURL][erc1155.TokenID] = database.NFTOwnership{
						Blockchain:       contractID.Blockchain,
						Network:          contractID.Network,
						ContractAddress:  contractID.Address,
						TokenID:          erc1155.TokenID,
						Balance:          balance,
						BlockNumber:      bunbig.FromMathBig(blockNumber.ToBigInt()),
						OwnerAddress:     transfer.To,
						TransactionHash:  transfer.Hash,
						TransactionIndex: uniqueID.TransactionIndex,
						BlockTimestamp:   &blockTime,
					}
				}
			}
			continue
		}

		// Transfer is ERC-721
		balance := contractIDToTokenIDToBalance[contractURL][transfer.TokenID]
		if _, ok := contractIDToTokenIDToOwnership[contractURL][transfer.TokenID]; !ok {
			contractIDToTokenIDToOwnership[contractURL][transfer.TokenID] = database.NFTOwnership{
				Blockchain:       contractID.Blockchain,
				Network:          contractID.Network,
				ContractAddress:  contractID.Address,
				TokenID:          transfer.TokenID,
				Balance:          balance,
				BlockNumber:      bunbig.FromMathBig(blockNumber.ToBigInt()),
				OwnerAddress:     transfer.To,
				TransactionHash:  transfer.Hash,
				TransactionIndex: uniqueID.TransactionIndex,
				BlockTimestamp:   &blockTime,
			}
		}
	}

	ownerships := make([]database.NFTOwnership, 0)
	for _, contract := range contracts {
		tokenIDs := contract.Query["token_ids"]
		strippedContractID := contract.StripQuery().String()

		contractOwnerships, ownershipsOk := contractIDToTokenIDToOwnership[strippedContractID]
		// Handle ERC-1155
		if len(tokenIDs) != 0 {
			// Append either existing ownership or empty ownership for each tokenID
			for _, tokenID := range tokenIDs {
				erc1155ownership, ok := contractOwnerships[tokenID]
				if !ownershipsOk || !ok {
					ownerships = append(ownerships, database.NewEmptyNFTOwnership(contract, tokenID, ownerID))
				}

				if ok {
					ownerships = append(ownerships, erc1155ownership)
				}
			}
		} else if ownershipsOk {
			for _, erc721ownership := range contractOwnerships {
				ownerships = append(ownerships, erc721ownership)
			}
		} else {
			ownerships = append(ownerships, database.NewEmptyNFTOwnership(contract, "0x0", ownerID))
		}

	}

	return ownerships, nil
}
