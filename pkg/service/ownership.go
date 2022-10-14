package service

import (
	"net/url"
	"time"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model/alchemy"
	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
	"github.com/authgear/authgear-server/pkg/util/clock"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

type OwnershipServiceNFTOwnershipMutator interface {
	InsertNFTOwnerships(ownerships []database.NFTOwnership) error
}

type OwnershipServiceAlchemyAPI interface {
	GetOwnerNFTs(ownerAddress string, contractIDs []authgearweb3.ContractID, pageKey string) (*alchemy.GetNFTsResponse, error)
	GetAssetTransfers(params web3.GetAssetTransferParams) (*alchemy.AssetTransferResult, error)
}

type OwnershipService struct {
	Clock               clock.Clock
	Config              config.Config
	AlchemyAPI          OwnershipServiceAlchemyAPI
	NFTCollectionQuery  query.NFTCollectionQuery
	NFTOwnershipQuery   query.NFTOwnershipQuery
	NFTOwnershipMutator OwnershipServiceNFTOwnershipMutator
}

func (h *OwnershipService) FetchAndInsertNFTOwnerships(ownerID authgearweb3.ContractID, contracts []authgearweb3.ContractID) ([]database.NFTOwnership, error) {
	pageKey := ""
	nftFetchCount := 0
	ownedNFTs := make([]alchemy.OwnedNFT, 0)
	contractIDsToEnquire := make([]authgearweb3.ContractID, 0)

	// Fetch user nfts until no extra page or has reached the page limit
	for ok := true; ok; ok = pageKey != "" && nftFetchCount <= h.Config.Server.MaxNFTPages {
		nfts, err := h.AlchemyAPI.GetOwnerNFTs(ownerID.Address.String(), contracts, pageKey)
		if err != nil {
			return []database.NFTOwnership{}, err
		}

		for _, ownedNFT := range nfts.OwnedNFTs {
			contractID, err := authgearweb3.NewContractID(ownerID.Blockchain, ownerID.Network, ownedNFT.Contract.Address, url.Values{})
			if err != nil {
				return []database.NFTOwnership{}, err
			}

			contractIDsToEnquire = append(contractIDsToEnquire, *contractID)
		}

		if nfts.PageKey != nil {
			pageKey = *nfts.PageKey
		}

		ownedNFTs = append(ownedNFTs, nfts.OwnedNFTs...)
		nftFetchCount++
	}

	nftTransfers := make([]alchemy.TokenTransfer, 0)
	if len(ownedNFTs) != 0 {
		pageKey = ""
		transferFetchCount := 0
		// Fetch transfers until no extra page or has reached the page limit
		for ok := true; ok; ok = pageKey != "" && transferFetchCount <= 5 {
			transfers, err := h.AlchemyAPI.GetAssetTransfers(web3.GetAssetTransferParams{
				ContractIDs: contractIDsToEnquire,
				ToAddress:   ownerID.Address,
				FromBlock:   "0x0",
				ToBlock:     "latest",
				PageKey:     pageKey,
				MaxCount:    1000,
				Order:       "desc",
			})
			if err != nil {
				return []database.NFTOwnership{}, err
			}
			nftTransfers = append(nftTransfers, transfers.Transfers...)
			transferFetchCount++
		}

	}

	ownerships, err := alchemy.MakeNFTOwnerships(ownerID, contracts, nftTransfers, ownedNFTs)
	if err != nil {
		return []database.NFTOwnership{}, err
	}

	// Insert ownerships
	err = h.NFTOwnershipMutator.InsertNFTOwnerships(ownerships)
	if err != nil {
		return []database.NFTOwnership{}, err
	}
	return ownerships, nil
}

func (h *OwnershipService) GetOwnerships(ownerID authgearweb3.ContractID, contracts []authgearweb3.ContractID) ([]database.NFTOwnership, error) {
	minimumFreshness := h.Clock.NowUTC()
	minimumFreshness = minimumFreshness.Add(-time.Duration(h.Config.Server.CacheTTL) * time.Second)

	// Query ownership from database
	ownershipQb := h.NFTOwnershipQuery.NewQueryBuilder()
	ownershipQb = ownershipQb.WithContracts(contracts).WithOwner(&ownerID).WithMinimumFreshness(minimumFreshness)
	ownerships, err := h.NFTOwnershipQuery.ExecuteQuery(ownershipQb)
	if err != nil {
		return []database.NFTOwnership{}, err
	}

	// Find out which contract to fetch
	contractsToFetch := make([]authgearweb3.ContractID, 0)
	contractIDToOwnerships := make(map[string][]database.NFTOwnership)
	for _, ownership := range ownerships {
		contractID := ownership.ContractID().String()

		if _, ok := contractIDToOwnerships[contractID]; ok {
			contractIDToOwnerships[contractID] = append(contractIDToOwnerships[contractID], ownership)
		} else {
			contractIDToOwnerships[contractID] = []database.NFTOwnership{ownership}
		}
	}

	for _, contract := range contracts {
		tokenIDs := contract.Query["token_ids"]

		strippedContractID := contract.StripQuery().String()

		ownerships, ok := contractIDToOwnerships[strippedContractID]
		if !ok || (len(tokenIDs) > 0 && len(ownerships) != len(tokenIDs)) {
			contractsToFetch = append(contractsToFetch, contract)
		}
	}

	// Fetch missing data from alchemy
	if len(contractsToFetch) != 0 {
		updatedOwnerships, err := h.FetchAndInsertNFTOwnerships(ownerID, contractsToFetch)
		if err != nil {
			return []database.NFTOwnership{}, err
		}

		for _, ownership := range updatedOwnerships {
			contractID := ownership.ContractID().String()

			if _, ok := contractIDToOwnerships[contractID]; ok {
				contractIDToOwnerships[contractID] = append(contractIDToOwnerships[contractID], ownership)
			} else {
				contractIDToOwnerships[contractID] = []database.NFTOwnership{ownership}
			}
		}
	}

	result := make([]database.NFTOwnership, 0)
	for _, contract := range contracts {
		contractID := contract.StripQuery().String()

		ownerships := contractIDToOwnerships[contractID]
		for _, ownership := range ownerships {
			if !ownership.IsEmpty() {
				result = append(result, ownership)
			}
		}

	}

	return result, nil
}
