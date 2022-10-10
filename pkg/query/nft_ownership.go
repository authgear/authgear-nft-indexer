package query

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun"
)

type NFTOwnershipQuery struct {
	Ctx     context.Context
	Session *bun.DB
}

type NFTOwnershipQueryBuilder struct {
	*bun.SelectQuery
}

func (b NFTOwnershipQueryBuilder) WithContracts(contracts []authgearweb3.ContractID) NFTOwnershipQueryBuilder {
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

	return NFTOwnershipQueryBuilder{
		b.Where("blockchain IN (?) AND network IN (?) AND contract_address IN (?)", bun.In(blockchains), bun.In(networks), bun.In(contractAddresses)),
	}
}

func (b NFTOwnershipQueryBuilder) WithTokenIDs(tokenIDs []string) NFTOwnershipQueryBuilder {
	if len(tokenIDs) == 0 {
		return b
	}

	return NFTOwnershipQueryBuilder{
		b.Where("token_id IN (?)", bun.In(tokenIDs)),
	}
}

func (b NFTOwnershipQueryBuilder) WithOwner(ownerID *authgearweb3.ContractID) NFTOwnershipQueryBuilder {
	if ownerID == nil {
		return b
	}
	return NFTOwnershipQueryBuilder{
		b.Where("blockchain = ? AND network = ? AND owner_address = ?", ownerID.Blockchain, ownerID.Network, ownerID.Address),
	}
}

func (q *NFTOwnershipQuery) NewQueryBuilder() NFTOwnershipQueryBuilder {
	return NFTOwnershipQueryBuilder{
		q.Session.NewSelect().Model((*database.NFTOwnership)(nil)),
	}
}

func (q *NFTOwnershipQuery) ExecuteQuery(qb NFTOwnershipQueryBuilder) ([]database.NFTOwnership, error) {
	nftOwnerships := make([]database.NFTOwnership, 0)

	query := qb.Order("token_id ASC")

	err := query.Scan(q.Ctx, &nftOwnerships)
	if err != nil {
		return []database.NFTOwnership{}, err
	}

	return nftOwnerships, nil
}
