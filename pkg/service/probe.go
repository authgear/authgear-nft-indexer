package service

import (
	"github.com/authgear/authgear-nft-indexer/pkg/model/alchemy"
	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

type ProbeServiceAlchemyAPI interface {
	GetOwnersForCollection(contractID authgearweb3.ContractID) (*alchemy.GetOwnersForCollectionResponse, error)
}

type ProbeServiceNFTCollectionProbeQuery interface {
	QueryCollectionProbeByContractID(contractID authgearweb3.ContractID) (*database.NFTCollectionProbe, error)
}

type ProbeServiceNFTCollectionProbeMutator interface {
	InsertNFTCollectionProbe(contractID authgearweb3.ContractID, isLargeCollection bool) (*database.NFTCollectionProbe, error)
}

type ProbeService struct {
	AlchemyAPI                ProbeServiceAlchemyAPI
	NFTCollectionProbeQuery   ProbeServiceNFTCollectionProbeQuery
	NFTCollectionProbeMutator ProbeServiceNFTCollectionProbeMutator
}

func (m *ProbeService) ProbeCollection(contractID authgearweb3.ContractID) (bool, error) {
	collectionProbe, err := m.NFTCollectionProbeQuery.QueryCollectionProbeByContractID(contractID)
	if err == nil && collectionProbe != nil {
		return collectionProbe.IsLargeCollection, nil
	}

	res, err := m.AlchemyAPI.GetOwnersForCollection(contractID)
	if err != nil {
		return false, err
	}

	dbProbe, err := m.NFTCollectionProbeMutator.InsertNFTCollectionProbe(contractID, res.PageKey != nil)
	if err != nil {
		return false, err
	}

	return dbProbe.IsLargeCollection, nil
}
