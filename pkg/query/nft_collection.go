package query

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/uptrace/bun"
)

type NFTCollectionQuery struct {
	Ctx     context.Context
	Session *bun.DB
}

type NFTCollectionQueryBuilder struct {
	*bun.SelectQuery
}

func (b NFTCollectionQueryBuilder) WithContracts(contracts []model.ContractID) NFTCollectionQueryBuilder {
	if len(contracts) == 0 {
		return b
	}

	blockchains := make([]string, len(contracts))
	networks := make([]string, len(contracts))
	contractAddresses := make([]string, len(contracts))

	for _, contract := range contracts {
		blockchains = append(blockchains, contract.Blockchain)
		networks = append(networks, contract.Network)
		contractAddresses = append(contractAddresses, contract.ContractAddress)
	}
	return NFTCollectionQueryBuilder{
		// We only support erc721 for now
		b.Where("blockchain IN (?) AND network IN (?) AND contract_address IN (?) AND type = ?", bun.In(blockchains), bun.In(networks), bun.In(contractAddresses), eth.NFTCollectionTypeERC721),
	}

}

func (q *NFTCollectionQuery) NewQueryBuilder() NFTCollectionQueryBuilder {
	return NFTCollectionQueryBuilder{
		q.Session.NewSelect().Model((*eth.NFTCollection)(nil)),
	}
}

func (q *NFTCollectionQuery) ExecuteQuery(qb NFTCollectionQueryBuilder) ([]eth.NFTCollection, error) {
	nftCollections := make([]eth.NFTCollection, 0)
	query := qb.Order("created_at DESC")
	err := query.Scan(q.Ctx, &nftCollections)
	if err != nil {
		return []eth.NFTCollection{}, err
	}

	return nftCollections, nil
}

func (q *NFTCollectionQuery) QueryAllNFTCollections() ([]eth.NFTCollection, error) {
	nftCollections := make([]eth.NFTCollection, 0)

	err := q.Session.NewSelect().Model(&nftCollections).Scan(q.Ctx)
	if err != nil {
		return nil, err
	}

	return nftCollections, nil
}
