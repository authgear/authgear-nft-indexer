package service

import (
	"fmt"
	"math/big"
	"net/url"

	"github.com/authgear/authgear-nft-indexer/pkg/model/alchemy"
	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

type MetadataServiceAlchemyAPI interface {
	GetContractMetadata(contractID authgearweb3.ContractID) (*alchemy.ContractMetadataResponse, error)
}

type MetadataServiceNFTCollectionMutator interface {
	InsertNFTCollection(contractID authgearweb3.ContractID, contractName string, tokenType database.NFTCollectionType, totalSupply *big.Int) (*database.NFTCollection, error)
}

type MetadataServiceRateLimiter interface {
	TakeToken(bucket ratelimit.Bucket) error
}

type MetadataService struct {
	AlchemyAPI           MetadataServiceAlchemyAPI
	NFTCollectionQuery   query.NFTCollectionQuery
	NFTCollectionMutator MetadataServiceNFTCollectionMutator
	RateLimiter          MetadataServiceRateLimiter
}

func (m *MetadataService) GetContractMetadata(appID string, contracts []authgearweb3.ContractID) ([]database.NFTCollection, error) {
	qb := m.NFTCollectionQuery.NewQueryBuilder()
	qb = qb.WithContracts(contracts)
	collections, err := m.NFTCollectionQuery.ExecuteQuery(qb)
	if err != nil {
		return []database.NFTCollection{}, err
	}

	contractIDToCollectionMap := make(map[string]*database.NFTCollection)
	for i, collection := range collections {
		contractID, err := collection.ContractID()
		if err != nil {
			return []database.NFTCollection{}, err
		}

		contractURL, err := contractID.URL()
		if err != nil {
			return []database.NFTCollection{}, err
		}

		contractIDToCollectionMap[contractURL.String()] = &collections[i]
	}

	res := make([]database.NFTCollection, 0, len(contracts))
	for _, contract := range contracts {
		contractID, err := authgearweb3.NewContractID(contract.Blockchain, contract.Network, contract.Address, url.Values{})
		if err != nil {
			return []database.NFTCollection{}, err
		}
		contractURL, err := contractID.URL()
		if err != nil {
			return []database.NFTCollection{}, err
		}

		err = m.RateLimiter.TakeToken(AntiSpamContractMetadataRequestBucket(appID))
		if err != nil {
			return []database.NFTCollection{}, err
		}

		// If exists, append to result, otherwise get from alchemy
		collection := contractIDToCollectionMap[contractURL.String()]
		if collection != nil {
			res = append(res, *collection)
			continue
		}

		contractMetadata, err := m.AlchemyAPI.GetContractMetadata(contract)
		if err != nil {
			return []database.NFTCollection{}, err
		}

		tokenType, err := database.ParseNFTCollectionType(contractMetadata.ContractMetadata.TokenType)
		if err != nil {
			return []database.NFTCollection{}, err
		}

		if contractMetadata.ContractMetadata.Name == "" {
			return []database.NFTCollection{}, fmt.Errorf("missing contract metadata")
		}

		totalSupply := new(big.Int)
		if contractMetadata.ContractMetadata.TotalSupply != "" {
			if _, ok := totalSupply.SetString(contractMetadata.ContractMetadata.TotalSupply, 10); !ok {
				return []database.NFTCollection{}, fmt.Errorf("failed to parse totalSupply")
			}
		}

		newCollection, err := m.NFTCollectionMutator.InsertNFTCollection(
			contract,
			contractMetadata.ContractMetadata.Name,
			tokenType,
			totalSupply,
		)

		if err != nil {
			return []database.NFTCollection{}, err
		}

		res = append(res, *newCollection)

	}

	return res, nil

}
