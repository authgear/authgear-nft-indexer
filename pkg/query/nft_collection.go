package query

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun"
)

type NFTCollectionQuery struct {
	Ctx     context.Context
	Session *bun.DB
}

type NFTCollectionQueryBuilder struct {
	*bun.SelectQuery
}

func (b NFTCollectionQueryBuilder) WithContracts(contracts []authgearweb3.ContractID) NFTCollectionQueryBuilder {
	if len(contracts) == 0 {
		return b
	}

	blockchains := make([]string, len(contracts))
	networks := make([]string, len(contracts))
	contractAddresses := make([]string, len(contracts))

	for _, contract := range contracts {
		blockchains = append(blockchains, contract.Blockchain)
		networks = append(networks, contract.Network)
		contractAddresses = append(contractAddresses, contract.Address)
	}
	return NFTCollectionQueryBuilder{
		b.Where("blockchain IN (?) AND network IN (?) AND contract_address IN (?)", bun.In(blockchains), bun.In(networks), bun.In(contractAddresses)),
	}

}

func (q *NFTCollectionQuery) NewQueryBuilder() NFTCollectionQueryBuilder {
	return NFTCollectionQueryBuilder{
		q.Session.NewSelect().Model((*database.NFTCollection)(nil)),
	}
}

func (q *NFTCollectionQuery) ExecuteQuery(qb NFTCollectionQueryBuilder) ([]database.NFTCollection, error) {
	nftCollections := make([]database.NFTCollection, 0)
	query := qb.Order("created_at DESC")
	err := query.Scan(q.Ctx, &nftCollections)
	if err != nil {
		return []database.NFTCollection{}, err
	}

	return nftCollections, nil
}

func (q *NFTCollectionQuery) QueryAllNFTCollections() ([]database.NFTCollection, error) {
	nftCollections := make([]database.NFTCollection, 0)

	err := q.Session.NewSelect().Model(&nftCollections).Scan(q.Ctx)
	if err != nil {
		return nil, err
	}

	return nftCollections, nil
}
