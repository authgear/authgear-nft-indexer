package query

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/uptrace/bun"
)

type NFTOwnerQuery struct {
	Ctx     context.Context
	Session *bun.DB
}

type NFTOwnerQueryBuilder struct {
	*bun.SelectQuery
}

func (b NFTOwnerQueryBuilder) WithContracts(contracts []model.ContractID) NFTOwnerQueryBuilder {

	blockchains := make([]string, len(contracts))
	networks := make([]string, len(contracts))
	contractAddresses := make([]string, len(contracts))

	for _, contract := range contracts {
		blockchains = append(blockchains, contract.Blockchain)
		networks = append(networks, contract.Network)
		contractAddresses = append(contractAddresses, contract.ContractAddress)
	}

	return NFTOwnerQueryBuilder{
		b.Where("blockchain IN (?) AND network IN (?) AND contract_address IN (?)", bun.In(blockchains), bun.In(networks), bun.In(contractAddresses)),
	}
}

func (b NFTOwnerQueryBuilder) WithOwnerAddress(ownerAddress string) NFTOwnerQueryBuilder {
	return NFTOwnerQueryBuilder{
		b.Where("owner_address = ?", ownerAddress),
	}
}

func (q *NFTOwnerQuery) NewQueryBuilder() NFTOwnerQueryBuilder {
	return NFTOwnerQueryBuilder{
		q.Session.NewSelect().Model((*ethmodel.NFTOwner)(nil)),
	}
}

func (q *NFTOwnerQuery) ExecuteQuery(qb NFTOwnerQueryBuilder) ([]ethmodel.NFTOwner, error) {
	nftOwners := make([]eth.NFTOwner, 0)

	query := qb.Order("token_id ASC")

	err := query.Scan(q.Ctx, &nftOwners)
	if err != nil {
		return []ethmodel.NFTOwner{}, err
	}

	return nftOwners, nil
}
