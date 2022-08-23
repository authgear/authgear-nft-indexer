package mutator

import (
	"context"
	"database/sql"

	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
)

type NFTCollectionMutator struct {
	Ctx     context.Context
	Session *bun.DB
}

func (q *NFTCollectionMutator) InsertNFTCollection(blockchain string, network string, contractName string, contractAddress string) (*eth.NFTCollection, error) {
	collection := &eth.NFTCollection{
		Blockchain:      blockchain,
		Network:         network,
		ContractAddress: contractAddress,
		Name:            contractName,
		FromBlockHeight: *bunbig.FromInt64(0),
	}

	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Query record ID if exists, if not, generate new id and insert
		tx.NewSelect().
			Model(collection).
			Where("blockchain = ? AND network = ? AND contract_address = ?", collection.Blockchain, collection.Network, collection.ContractAddress).
			Limit(1).
			Scan(ctx) //nolint:errcheck

		_, err := tx.NewInsert().
			Model(collection).
			On("CONFLICT (blockchain, network, contract_address) DO NOTHING").
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
