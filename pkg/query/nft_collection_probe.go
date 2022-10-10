package query

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun"
)

type NFTCollectionProbeQuery struct {
	Ctx     context.Context
	Session *bun.DB
}

func (q *NFTCollectionProbeQuery) QueryCollectionProbeByContractID(contractID authgearweb3.ContractID) (*database.NFTCollectionProbe, error) {
	probe := new(database.NFTCollectionProbe)

	err := q.Session.NewSelect().Model(probe).Where(
		"blockchain = ? AND network = ? AND contract_address = ?", contractID.Blockchain, contractID.Network, contractID.Address,
	).Scan(q.Ctx)
	if err != nil {
		return nil, err
	}

	return probe, nil
}
