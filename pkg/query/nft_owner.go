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

func (b NFTOwnerQueryBuilder) WithBlockchainNetwork(blockchainNetwork model.BlockchainNetwork) NFTOwnerQueryBuilder {
	return NFTOwnerQueryBuilder{
		b.Where("blockchain = ? AND network = ?", blockchainNetwork.Blockchain, blockchainNetwork.Network),
	}
}

func (b NFTOwnerQueryBuilder) WithContractAddresses(contractAddresses []string) NFTOwnerQueryBuilder {
	return NFTOwnerQueryBuilder{
		b.Where("contract_address IN (?)", bun.In(contractAddresses)),
	}
}

func (b NFTOwnerQueryBuilder) WithOwnerAddresses(ownerAddresses []string) NFTOwnerQueryBuilder {
	return NFTOwnerQueryBuilder{
		b.Where("owner_address IN (?)", bun.In(ownerAddresses)),
	}
}

func (q *NFTOwnerQuery) NewQueryBuilder() NFTOwnerQueryBuilder {
	return NFTOwnerQueryBuilder{
		q.Session.NewSelect().Model((*ethmodel.NFTOwner)(nil)),
	}
}

func (q *NFTOwnerQuery) ExecuteQuery(qb NFTOwnerQueryBuilder, limit int, offset int) (model.Paginated[ethmodel.NFTOwner], error) {
	nftOwners := make([]eth.NFTOwner, 0)

	query := qb.Order("token_id ASC")

	totalCount, err := query.Count(q.Ctx)

	if err != nil {
		return model.Paginated[eth.NFTOwner]{
			Items:      []ethmodel.NFTOwner{},
			TotalCount: 0,
		}, err
	}

	err = query.Limit(limit).Offset(offset).Scan(q.Ctx, &nftOwners)
	if err != nil {
		return model.Paginated[eth.NFTOwner]{
			Items:      []ethmodel.NFTOwner{},
			TotalCount: 0,
		}, err
	}

	return model.Paginated[eth.NFTOwner]{
		Items:      nftOwners,
		TotalCount: totalCount,
	}, nil
}
