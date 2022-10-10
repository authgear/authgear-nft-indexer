package mutator

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun"
)

type NFTCollectionProbeMutator struct {
	Ctx     context.Context
	Session *bun.DB
}

func (q *NFTCollectionProbeMutator) InsertNFTCollectionProbe(contractID authgearweb3.ContractID, isLargeCollection bool) (*database.NFTCollectionProbe, error) {
	probe := &database.NFTCollectionProbe{
		Blockchain:        contractID.Blockchain,
		Network:           contractID.Network,
		ContractAddress:   contractID.Address,
		IsLargeCollection: isLargeCollection,
	}

	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(probe).
			On("CONFLICT (blockchain, network, contract_address) DO NOTHING").
			Returning("*").
			Exec(ctx)
		return err
	})

	if err != nil {
		return nil, err
	}

	return probe, nil
}
