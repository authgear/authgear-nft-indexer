package mutator

import (
	"context"
	"database/sql"

	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
)

type NFTCollectionMutator struct {
	Ctx     context.Context
	Session *bun.DB
}

func (q *NFTCollectionMutator) InsertNFTCollection(blockchainNetwork model.BlockchainNetwork, contractName string, contractAddress string) (*eth.NFTCollection, error) {
	collection := &eth.NFTCollection{
		Blockchain:        blockchainNetwork.Blockchain,
		Network:           blockchainNetwork.Network,
		ContractAddress:   contractAddress,
		Name:              contractName,
		SyncedBlockHeight: *bunbig.FromInt64(0),
	}

	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(collection).
			Returning("*").
			Exec(ctx)
		return err
	})

	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (q *NFTCollectionMutator) DeleteNFTCollection(id string) (*eth.NFTCollection, error) {

	collection := &eth.NFTCollection{}

	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewDelete().
			Model(collection).
			Where("id = ?", id).
			Returning("*").
			Exec(ctx)

		if row, err := res.RowsAffected(); err != nil || row == 0 {
			return sql.ErrNoRows
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return collection, nil
}
