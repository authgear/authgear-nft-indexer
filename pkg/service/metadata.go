package service

import (
	"math/big"
	"time"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model/alchemy"
	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/clock"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

type MetadataServiceAlchemyAPI interface {
	GetContractMetadata(contractID authgearweb3.ContractID) (*alchemy.ContractMetadataResponse, error)
}

type MetadataServiceNFTCollectionMutator interface {
	InsertNFTCollection(contractID authgearweb3.ContractID, contractName string, tokenType database.NFTCollectionType, totalSupply *big.Int) (*database.NFTCollection, error)
}

type MetadataService struct {
	Clock                clock.Clock
	Config               config.Config
	AlchemyAPI           MetadataServiceAlchemyAPI
	NFTCollectionQuery   query.NFTCollectionQuery
	NFTCollectionMutator MetadataServiceNFTCollectionMutator
}

func (m *MetadataService) GetContractMetadata(contracts []authgearweb3.ContractID) ([]database.NFTCollection, error) {
	minimumFreshness := m.Clock.NowUTC()
	minimumFreshness = minimumFreshness.Add(-time.Duration(m.Config.Server.CollectionCacheTTL) * time.Second)

	qb := m.NFTCollectionQuery.NewQueryBuilder()
	qb = qb.WithContracts(contracts).WithMinimumFreshness(minimumFreshness)
	collections, err := m.NFTCollectionQuery.ExecuteQuery(qb)
	if err != nil {
		return nil, err
	}

	contractIDToCollectionMap := make(map[string]*database.NFTCollection)
	for i, collection := range collections {
		contractID := collection.ContractID().String()

		contractIDToCollectionMap[contractID] = &collections[i]
	}

	res := make([]database.NFTCollection, 0, len(contracts))
	for _, contract := range contracts {
		strippedContractID := contract.StripQuery().String()
		// If exists, append to result, otherwise get from alchemy
		collection := contractIDToCollectionMap[strippedContractID]
		if collection != nil {
			res = append(res, *collection)
			continue
		}

		contractMetadata, err := m.AlchemyAPI.GetContractMetadata(contract)
		if err != nil {
			return nil, err
		}

		tokenType, err := database.ParseNFTCollectionType(contractMetadata.ContractMetadata.TokenType)
		if err != nil {
			return nil, ErrBadNFTCollection.NewWithDetails("unable to parse token type", apierrors.Details{"tokenType": contractMetadata.ContractMetadata.TokenType})
		}

		if contractMetadata.ContractMetadata.Name == "" {
			return nil, ErrBadNFTCollection.New("missing contract metadata")
		}

		totalSupply := new(big.Int)
		if contractMetadata.ContractMetadata.TotalSupply != "" {
			if _, ok := totalSupply.SetString(contractMetadata.ContractMetadata.TotalSupply, 10); !ok {
				return nil, ErrBadNFTCollection.NewWithDetails("failed to parse total supply", apierrors.Details{"totalSupply": contractMetadata.ContractMetadata.TotalSupply})
			}
		}

		newCollection, err := m.NFTCollectionMutator.InsertNFTCollection(
			contract,
			contractMetadata.ContractMetadata.Name,
			tokenType,
			totalSupply,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, *newCollection)

	}

	return res, nil

}
