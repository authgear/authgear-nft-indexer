package query

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun"
)

type NFTOwnerQuery struct {
	Ctx     context.Context
	Session *bun.DB
}

func (q *NFTOwnerQuery) QueryOwner(ownerID authgearweb3.ContractID) (*database.NFTOwner, error) {
	owner := new(database.NFTOwner)

	err := q.Session.NewSelect().
		Model(owner).
		Where("blockchain = ? AND network = ? AND address = ?", ownerID.Blockchain, ownerID.Network, ownerID.Address).
		Scan(q.Ctx)
	if err != nil {
		return nil, err
	}

	return owner, nil
}
